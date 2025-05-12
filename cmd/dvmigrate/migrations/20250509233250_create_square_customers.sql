-- +goose Up
-- +goose StatementBegin
CREATE TABLE square_customers (
  square_merchant_id TEXT NOT NULL REFERENCES square_merchants (square_merchant_id),
  square_customer_id TEXT PRIMARY KEY,
  given_name TEXT NOT NULL,
  family_name TEXT NOT NULL,
  email_address TEXT NOT NULL,
  phone_number TEXT NOT NULL,
  customer_uuid UUID,
  created_at TIMESTAMP NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE square_customers;

-- +goose StatementEnd
