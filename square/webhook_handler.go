package square

import (
	"fmt"
	"net/http"
)

const squareSignatureHeaderKey = "x-square-hmacsha256-signature"

type WebhookEvent struct{}

func HandleWebhookEvents() http.HandlerFunc {
	// TODO: pull from db
	signatureKey := "test"
	// TODO: pull from config
	webhookUrl := "https://example.com/webhook_events"

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var bodyBytes []byte
			_, err := r.Body.Read(bodyBytes)
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				fmt.Fprint(w, "unable to read request body") //nolint:errcheck
			}

			requestBody := string(bodyBytes)
			requestSignature := r.Header.Get(squareSignatureHeaderKey)

			if err := ValidateWebhookEvent(webhookUrl, requestBody, requestSignature, signatureKey); err != nil {
				// we write success because we don't want square to retry
				// and we also don't want people to brute force signatures
				w.WriteHeader(http.StatusOK)
			}

			// TODO: check event type
			// TODO: write event to database

			w.WriteHeader(http.StatusOK)
		},
	)
}
