package internalhttp //nolint

import (
	"net/http"
)

type Handler struct {
	logger Logger
	app    Application
}

func NewHandler(logger Logger, app Application) *Handler {
	return &Handler{logger: logger, app: app}
}

func (h *Handler) helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
//nolint