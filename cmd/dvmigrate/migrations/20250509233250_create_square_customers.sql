-- +goose Up
-- +goose StatementBegin
CREATE TABLE square_customers (
  square_merchant_id TEXT REFERENCES square_merchants (square_merchant_id),
  square_customer_id TEXT PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL,
  customer_uuid UUID,
  created_at TIMESTAMP NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE square_customers;

-- +goose StatementEnd
