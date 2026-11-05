package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"nexo/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// writeError traduce errores del dominio a códigos HTTP.
func writeError(w http.ResponseWriter, err error) {
	var ve domain.ValidationError
	status := http.StatusInternalServerError
	msg := "error interno del servidor"
	switch {
	case errors.As(err, &ve):
		status, msg = http.StatusUnprocessableEntity, ve.Msg
	case errors.Is(err, domain.ErrNotFound):
		status, msg = http.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrInvalidCredentials), errors.Is(err, domain.ErrUnauthorized):
		status, msg = http.StatusUnauthorized, err.Error()
	case errors.Is(err, domain.ErrInactive), errors.Is(err, domain.ErrForbidden):
		status, msg = http.StatusForbidden, err.Error()
	case errors.Is(err, domain.ErrEmailTaken):
		status, msg = http.StatusConflict, err.Error()
	default:
		log.Printf("error: %v", err)
	}
	writeJSON(w, status, map[string]string{"error": msg})
}

func decode(r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return domain.ValidationError{Msg: "el cuerpo de la petición no es válido"}
	}
	return nil
}
