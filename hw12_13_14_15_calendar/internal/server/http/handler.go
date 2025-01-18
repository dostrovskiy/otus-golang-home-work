package internalhttp //nolint

import "net/http"

type Handler struct {
	logger Logger
}

func NewHandler(logger Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) helloHandler(w http.ResponseWriter, r *http.Request) {
	_ = r
	w.WriteHeader(http.StatusNoContent)
}
//nolint