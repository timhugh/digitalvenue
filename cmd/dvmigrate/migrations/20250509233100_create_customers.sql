-- +goose Up
-- +goose StatementBegin
CREATE TABLE customers (
  tenant_uuid UUID NOT NULL REFERENCES tenants (tenant_uuid),
  customer_uuid UUID PRIMARY KEY,
  given_name TEXT NOT NULL,
  family_name TEXT NOT NULL,
  email_address TEXT NOT NULL,
  phone_number TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE customers;

-- +goose StatementEnd
