package webhooks

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
)

func Validate(body string, notification_url string, signature_key string, signature string) error {
	payload := new(bytes.Buffer)
	err := json.Compact(payload, []byte(body))
	if err != nil {
		return err
	}

	appended := append([]byte(notification_url), payload.Bytes()...)

	goodSignature := generateSignature(signature_key, appended)
	if goodSignature != signature {
		log.Debug().
			Str("expectedSignature", goodSignature).
			Str("actualSignature", signature).
			Msg("Signature mismatch")
		return fmt.Errorf("signature mismatch")
	}

	return nil
}

func generateSignature(key string, payload []byte) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write(payload)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
