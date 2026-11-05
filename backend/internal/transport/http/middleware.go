package http

import (
	"context"
	"net/http"
	"strings"

	"nexo/internal/domain"
)

type ctxKey struct{}

// ProfileLoader es lo único que el middleware necesita del caso de uso de auth.
type ProfileLoader interface {
	Profile(ctx context.Context, userID int64) (*domain.Profile, error)
}

type Middleware struct {
	tokens   domain.TokenManager
	profiles ProfileLoader
}

func NewMiddleware(t domain.TokenManager, p ProfileLoader) *Middleware {
	return &Middleware{tokens: t, profiles: p}
}

func profileFrom(ctx context.Context) *domain.Profile {
	p, _ := ctx.Value(ctxKey{}).(*domain.Profile)
	return p
}

// Authenticate valida el JWT y deja el perfil (con permisos vigentes) en el contexto.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			writeError(w, domain.ErrUnauthorized)
			return
		}
		id, err := m.tokens.Parse(token)
		if err != nil {
			writeError(w, domain.ErrUnauthorized)
			return
		}
		p, err := m.profiles.Profile(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, p)))
	})
}

// Require exige un permiso concreto. Se encadena después de Authenticate.
func (m *Middleware) Require(code string, next http.HandlerFunc) http.Handler {
	return m.RequireAny([]string{code}, next)
}

// RequireAny deja pasar si el perfil tiene al menos uno de los permisos indicados.
func (m *Middleware) RequireAny(codes []string, next http.HandlerFunc) http.Handler {
	return m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := profileFrom(r.Context())
		if p != nil {
			for _, c := range codes {
				if p.Has(c) {
					next(w, r)
					return
				}
			}
		}
		writeError(w, domain.ErrForbidden)
	}))
}

func CORS(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
