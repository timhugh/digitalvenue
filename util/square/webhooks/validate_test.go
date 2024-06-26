package webhooks

import (
	"github.com/matryer/is"
	"testing"
)

const (
	signatureKey    = "asdf1234"
	requestBody     = `{"hello":"world"}`
	notificationURL = "https://example.com/webhook"
	goodSignature   = "2kRE5qRU2tR+tBGlDwMEw2avJ7QM4ikPYD/PJ3bd9Og="
)

func TestValidate_noErrOnGoodSignature(t *testing.T) {
	is := is.New(t)
	err := Validate(requestBody, notificationURL, signatureKey, goodSignature)
	is.NoErr(err)
}

func TestValidate_errOnBadSignature(t *testing.T) {
	is := is.New(t)
	err := Validate(requestBody, notificationURL, signatureKey, "bad signature")
	is.True(err != nil)
}
