package config

import (
	"os"
	"time"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	JWTTTL        time.Duration
	JWTRememberTTL time.Duration
	AllowedOrigin string
	AdminName     string
	AdminEmail    string
	AdminPassword string
}

func Load() Config {
	ttl, err := time.ParseDuration(env("JWT_TTL", "8h"))
	if err != nil {
		ttl = 8 * time.Hour
	}
	remember, err := time.ParseDuration(env("JWT_REMEMBER_TTL", "720h"))
	if err != nil {
		remember = 720 * time.Hour
	}
	return Config{
		Port:          env("PORT", "8080"),
		DatabaseURL:   env("DATABASE_URL", "postgres://postgres:123456789@localhost:5432/nexo?sslmode=disable"),
		JWTSecret:     env("JWT_SECRET", "cambia-este-secreto-en-produccion"),
		JWTTTL:        ttl,
		JWTRememberTTL: remember,
		AllowedOrigin: env("ALLOWED_ORIGIN", "http://localhost:5173"),
		AdminName:     env("ADMIN_NAME", "Administrador"),
		AdminEmail:    env("ADMIN_EMAIL", "admin@nexo.local"),
		AdminPassword: env("ADMIN_PASSWORD", "Admin12345"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
