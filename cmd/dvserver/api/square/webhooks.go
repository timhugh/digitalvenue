package square

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/timhugh/digitalvenue/common"
	"github.com/timhugh/digitalvenue/db/queries"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/util/events"
)

const squareSignatureHeaderKey = "x-square-hmacsha256-signature"

type webhookEventDB interface {
	GetSquareMerchant(context.Context, string) (queries.SquareMerchant, error)
	GetSquarePayment(context.Context, string) (queries.SquarePayment, error)
	GetSquareCustomer(context.Context, string) (queries.SquareCustomer, error)
	CreateSquareCustomer(context.Context, queries.CreateSquareCustomerParams) (queries.SquareCustomer, error)
	CreateSquarePayment(context.Context, queries.CreateSquarePaymentParams) (queries.SquarePayment, error)
}

type webhookEvent struct {
	SquareMerchantID string `json:"merchant_id"`
	Type             string `json:"type"`

	Data struct {
		Type string `json:"type"`
		ID   string `json:"id"`

		Object struct {
			Payment struct {
				LocationID string `json:"location_id"`
				OrderID    string `json:"order_id"`
			} `json:"payment"`

			Customer struct {
				GivenName    string `json:"given_name"`
				FamilyName   string `json:"family_name"`
				EmailAddress string `json:"email_address"`
				PhoneNumber  string `json:"phone_number"`
			} `json:"customer"`
		} `json:"object"`
	} `json:"data"`
}

type processError struct {
	error
	StatusCode int
}

func HandleWebhookEvents(log common.Logger, db webhookEventDB, bus *events.Bus) http.Handler {
	webhookUrl := os.Getenv("DV_SQUARE_WEBHOOK_URL")

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				log.Error("unable to read request body: %w", err)
				return
			}

			requestBody := string(bodyBytes)
			requestSignature := r.Header.Get(squareSignatureHeaderKey)

			var event webhookEvent
			if err := json.Unmarshal(bodyBytes, &event); err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				log.Error("unable to parse request body json: %w", err)
				return
			}

			merchant, err := db.GetSquareMerchant(r.Context(), event.SquareMerchantID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					// we write success because we don't want square to retry
					// and we also don't want people to brute force requests
					//
					// TODO: technically we should still attempt signature validation
					// for timing attacks, but we're not going to worry about that right now
					w.WriteHeader(http.StatusOK)
					log.Warn("received webhook event for unknown merchant: %s", event.SquareMerchantID)
					return
				}

				// otherwise this might be a db connection issue or whatever else, so we should let square retry
				w.WriteHeader(http.StatusInternalServerError)
				log.Error("failed to fetch square merchant '%s' from db: %w", event.SquareMerchantID, err)
				return
			}

			if err := square.ValidateWebhookEvent(webhookUrl, requestBody, requestSignature, merchant.SquareWebhookSignatureKey); err != nil {
				// we write success because we don't want square to retry
				// and we also don't want people to brute force requests
				w.WriteHeader(http.StatusOK)
				return
			}

			var processErr *processError
			switch event.Type {
			case "payment.created":
				processErr = processPaymentCreated(r.Context(), db, bus, event)
			case "customer.created":
				processErr = processCustomerCreated(r.Context(), db, bus, event)
			}

			if processErr != nil {
				w.WriteHeader(processErr.StatusCode)
				log.Error("failed to process webhook event: %w", processErr)
				return
			}

			w.WriteHeader(http.StatusOK)
		},
	)
}

func processCustomerCreated(ctx context.Context, db webhookEventDB, bus *events.Bus, event webhookEvent) *processError {
	_, existingCustomerErr := db.GetSquareCustomer(ctx, event.Data.ID)

	if existingCustomerErr != nil && !errors.Is(existingCustomerErr, sql.ErrNoRows) {
		return &processError{existingCustomerErr, http.StatusInternalServerError}
	}

	if existingCustomerErr != nil && errors.Is(existingCustomerErr, sql.ErrNoRows) {
		createParams := queries.CreateSquareCustomerParams{
			SquareMerchantID: event.SquareMerchantID,
			SquareCustomerID: event.Data.ID,
			GivenName:        event.Data.Object.Customer.GivenName,
			FamilyName:       event.Data.Object.Customer.FamilyName,
			EmailAddress:     event.Data.Object.Customer.EmailAddress,
			PhoneNumber:      event.Data.Object.Customer.PhoneNumber,
		}

		if _, err := db.CreateSquareCustomer(ctx, createParams); err != nil {
			return &processError{err, http.StatusInternalServerError}
		}
	}

	if errs := events.Emit(bus, square.CustomerCreated{
		SquareMerchantID: event.SquareMerchantID,
		SquareCustomerID: event.Data.ID,
	}); len(errs) > 0 {
		err := fmt.Errorf("failed to emit square.CustomerCreated event: %v", errs)
		return &processError{err, http.StatusInternalServerError}
	}

	return nil
}

func processPaymentCreated(ctx context.Context, db webhookEventDB, bus *events.Bus, event webhookEvent) *processError {
	_, existingPaymentErr := db.GetSquarePayment(ctx, event.Data.ID)

	if existingPaymentErr != nil && !errors.Is(existingPaymentErr, sql.ErrNoRows) {
		return &processError{existingPaymentErr, http.StatusInternalServerError}
	}

	if existingPaymentErr != nil && errors.Is(existingPaymentErr, sql.ErrNoRows) {
		createParams := queries.CreateSquarePaymentParams{
			SquarePaymentID:  event.Data.ID,
			SquareMerchantID: event.SquareMerchantID,
		}

		if _, err := db.CreateSquarePayment(ctx, createParams); err != nil {
			return &processError{err, http.StatusInternalServerError}
		}
	}

	if errs := events.Emit(bus, square.PaymentCreated{
		SquareMerchantID: event.SquareMerchantID,
		SquarePaymentID:  event.Data.ID,
	}); len(errs) > 0 {
		err := fmt.Errorf("failed to emit square.PaymentCreated event: %v", errs)
		return &processError{err, http.StatusInternalServerError}
	}

	return nil
}
