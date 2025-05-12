-- +goose Up
-- +goose StatementBegin
CREATE TABLE square_merchants (
  tenant_uuid UUID NOT NULL REFERENCES tenants (tenant_uuid),
  square_merchant_id TEXT PRIMARY KEY,
  square_api_token TEXT NOT NULL,
  square_webhook_signature_key TEXT NOT NULL,
  ticketable_categories JSONB NOT NULL DEFAULT '[]',
  location_ids JSONB NOT NULL DEFAULT '[]',
  created_at TIMESTAMP NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE square_merchants;

-- +goose StatementEnd
