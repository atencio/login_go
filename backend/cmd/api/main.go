package main

import (
	"context"
	"errors"
	"log"
	nethttp "net/http"
	"os/signal"
	"syscall"
	"time"

	"nexo/internal/config"
	"nexo/internal/repository/postgres"
	"nexo/internal/security"
	httpx "nexo/internal/transport/http"
	"nexo/internal/usecase"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := postgres.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("base de datos: %v", err)
	}
	defer db.Close()
	if err := postgres.Migrate(ctx, db); err != nil {
		log.Fatalf("%v", err)
	}

	// Composición: aquí (y solo aquí) se conectan las implementaciones concretas.
	userRepo, roleRepo, reportRepo := postgres.NewUserRepo(db), postgres.NewRoleRepo(db), postgres.NewReportRepo(db)
	positionRepo, employeeRepo := postgres.NewPositionRepo(db), postgres.NewEmployeeRepo(db)
	equipmentRepo, assignmentRepo := postgres.NewEquipmentRepo(db), postgres.NewAssignmentRepo(db)
	hasher := security.NewBcryptHasher()
	tokens := security.NewJWTManager(cfg.JWTSecret, cfg.JWTTTL, cfg.JWTRememberTTL)

	authSvc := usecase.NewAuthService(userRepo, roleRepo, hasher, tokens)
	userSvc := usecase.NewUserService(userRepo, roleRepo, hasher)
	roleSvc := usecase.NewRoleService(roleRepo)
	reportSvc := usecase.NewReportService(reportRepo)
	positionSvc := usecase.NewPositionService(positionRepo)
	employeeSvc := usecase.NewEmployeeService(employeeRepo)
	equipmentSvc := usecase.NewEquipmentService(equipmentRepo)
	assignmentSvc := usecase.NewAssignmentService(assignmentRepo, equipmentRepo, employeeRepo)

	if created, err := userSvc.EnsureAdmin(ctx, cfg.AdminName, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Fatalf("admin inicial: %v", err)
	} else if created {
		log.Printf("usuario administrador creado: %s", cfg.AdminEmail)
	}

	handler := httpx.NewHandler(authSvc, userSvc, roleSvc, reportSvc)
	inventory := httpx.NewInventoryHandler(positionSvc, employeeSvc, equipmentSvc, assignmentSvc)
	router := httpx.NewRouter(handler, inventory, httpx.NewMiddleware(tokens, authSvc), cfg.AllowedOrigin)

	srv := &nethttp.Server{Addr: ":" + cfg.Port, Handler: router, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdown)
	}()

	log.Printf("API escuchando en :%s", cfg.Port)
	if err := srv.ListenAndServe(); !errors.Is(err, nethttp.ErrServerClosed) {
		log.Fatal(err)
	}
}
