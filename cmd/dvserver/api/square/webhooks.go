package square

import (
	"io"
	"net/http"
)

func HandleWebhookEvents() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "OK") //nolint:errcheck
		},
	)
}
