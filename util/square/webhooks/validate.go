package webhooks

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func Validate(body string, notificationURL string, signatureKey string, signature string) error {
	payload := new(bytes.Buffer)
	err := json.Compact(payload, []byte(body))
	if err != nil {
		return errors.Wrap(err, "failed to compact json")
	}

	appended := append([]byte(notificationURL), payload.Bytes()...)

	goodSignature := generateSignature(signatureKey, appended)
	if goodSignature != signature {
		log.Debug().
			Str("expectedSignature", goodSignature).
			Str("actualSignature", signature).
			Msg("Signature mismatch")
		return errors.New("signature mismatch")
	}

	return nil
}

func generateSignature(key string, payload []byte) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write(payload)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
