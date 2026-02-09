package httpapi

import (
	"errors"
	"net/http"

	"toolmesh/internal/domain"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrBadRequest):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrToolFailed):
		writeJSON(w, http.StatusBadGateway, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrModelUnavailable):
		writeJSON(w, http.StatusServiceUnavailable, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrNoContextFound):
		writeJSON(w, http.StatusOK, errorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}
