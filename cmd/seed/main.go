package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/timhugh/digitalvenue/db"
	dvDB "github.com/timhugh/digitalvenue/db"
	"github.com/timhugh/digitalvenue/db/queries"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	dbURI := os.Getenv("DV_DATABASE_URL")
	dbClient, err := db.New(dbURI)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}
	db := queries.New(dbClient)

	_, err = dvDB.CreateUser(ctx, db, queries.CreateUserParams{
		UserUuid: uuid.MustParse(os.Getenv("DV_USER_UUID")),
		Username: os.Getenv("DV_USER_NAME"),
		Name:     os.Getenv("DV_USER_NAME"),
	}, os.Getenv("DV_USER_PASSWORD"))
	if err != nil {
		panic(fmt.Errorf("failed to create user: %w", err))
	}

	tenant, err := db.CreateTenant(ctx, queries.CreateTenantParams{
		TenantUuid: uuid.MustParse(os.Getenv("DV_TENANT_UUID")),
		Name:       os.Getenv("DV_TENANT_NAME"),
	})
	if err != nil {
		panic(fmt.Errorf("failed to create tenant: %w", err))
	}
	_, err = db.CreateSquareMerchant(ctx, queries.CreateSquareMerchantParams{
		TenantUuid:                tenant.TenantUuid,
		SquareMerchantID:          os.Getenv("DV_SQUARE_MERCHANT_ID"),
		SquareApiToken:            os.Getenv("DV_SQUARE_ACCESS_TOKEN"),
		SquareWebhookSignatureKey: os.Getenv("DV_SQUARE_WEBHOOK_SIGNATURE_KEY"),
		TicketableCategories:      []string{},
		LocationIds:               []string{},
	})
	if err != nil {
		panic(fmt.Errorf("failed to create square merchant: %w", err))
	}

}
