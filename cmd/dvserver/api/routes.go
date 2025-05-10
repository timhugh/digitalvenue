package api

import (
	"net/http"

	"github.com/timhugh/digitalvenue/util/events"

	"github.com/timhugh/digitalvenue/cmd/dvserver/api/square"
)

func addRoutes(
	mux *http.ServeMux,
	_ *events.Bus,
) {
	mux.Handle("GET /health", handleHealthGet())
	mux.Handle("POST /square/webhook_events", square.HandleWebhookEvents())
}
