package square_test

import (
	"errors"
	"testing"

	"github.com/timhugh/digitalvenue/square"
)

func TestWebhookSignatureValidator(t *testing.T) {
	const webhookUrl = "https://example.com/webhook"
	const requestBody = "{\"hello\":\"world\"}"
	const signatureKey = "asdf1234"

	t.Run("valid signature", func(t *testing.T) {
		const validSignature = "2kRE5qRU2tR+tBGlDwMEw2avJ7QM4ikPYD/PJ3bd9Og="

		err := square.ValidateWebhookEvent(webhookUrl, requestBody, validSignature, signatureKey)
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}
	})

	t.Run("invalid signature", func(t *testing.T) {
		const invalidSignature = "bad signature"

		err := square.ValidateWebhookEvent(webhookUrl, requestBody, invalidSignature, signatureKey)
		if err == nil {
			t.Fatalf("Expected %s, got nil", square.ErrInvalidWebhookSignature)
		}
		if !errors.Is(err, square.ErrInvalidWebhookSignature) {
			t.Errorf("Expected %s, got %s", square.ErrInvalidWebhookSignature, err)
		}
	})
}
