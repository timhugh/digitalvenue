package api

import (
	"net/http"

	"github.com/timhugh/digitalvenue/square"
)

func addRoutes(
	mux *http.ServeMux,
) {
	mux.Handle("GET /health", handleHealthGet())
	mux.Handle("POST /square/webhook_events", square.HandleWebhookEvents())
}
