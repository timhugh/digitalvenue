package api

import (
	"net/http"

	"github.com/timhugh/digitalvenue/common"
	"github.com/timhugh/digitalvenue/db/queries"
	"github.com/timhugh/digitalvenue/util/events"

	"github.com/timhugh/digitalvenue/cmd/dvserver/api/square"
)

func addRoutes(
	mux *http.ServeMux,
	log common.Logger,
	bus *events.Bus,
	queries *queries.Queries,
) {
	mux.Handle("GET /health", handleHealthGet())
	mux.Handle("POST /square/webhook_events", square.HandleWebhookEvents(log, queries, bus))
}
