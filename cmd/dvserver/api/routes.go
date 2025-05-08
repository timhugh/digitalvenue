package api

import (
	"io"
	"net/http"

	"github.com/timhugh/digitalvenue/util/events"

	"github.com/timhugh/digitalvenue/square"
)

func testRoute(bus *events.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events.Emit(bus, "test")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK") //nolint:errcheck
	}
}

func addRoutes(
	mux *http.ServeMux,
	bus *events.Bus,
) {
	mux.Handle("GET /health", handleHealthGet())
	mux.Handle("POST /square/webhook_events", square.HandleWebhookEvents())
	mux.Handle("GET /test", testRoute(bus))
}
