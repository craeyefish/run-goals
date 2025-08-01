package handlers

import (
	"log"
	"net/http"
	"run-goals/controllers"
	"strings"
)

type SupportHandler struct {
	l                 *log.Logger
	supportController *controllers.SupportController
}

func NewSupportHandler(
	l *log.Logger,
	supportController *controllers.SupportController,
) *SupportHandler {
	return &SupportHandler{
		l:                 l,
		supportController: supportController,
	}
}

func (h *SupportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers for support endpoints (since they may be called from different domains)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/support/")

	switch {
	case strings.HasPrefix(path, "delete-account/"):
		h.supportController.DeleteUserAccount(w, r)
	default:
		h.l.Printf("Unsupported support endpoint: %s", path)
		http.Error(w, "Endpoint not found", http.StatusNotFound)
	}
}
