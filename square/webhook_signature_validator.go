package square

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

var ErrInvalidWebhookSignature = errors.New("invalid webhook signature")

func ValidateWebhookEvent(webhookUrl, requestBody, signature, signatureKey string) error {
	payload := new(bytes.Buffer)
	err := json.Compact(payload, []byte(requestBody))
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	appended := append([]byte(webhookUrl), payload.Bytes()...)

	hash := hmac.New(sha256.New, []byte(signatureKey))
	hash.Write(appended)
	goodSignature := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	if goodSignature != signature {
		return ErrInvalidWebhookSignature
	}

	return nil
}
