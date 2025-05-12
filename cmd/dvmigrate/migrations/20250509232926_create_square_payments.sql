-- +goose Up
-- +goose StatementBegin
CREATE TABLE square_payments (
  square_merchant_id TEXT NOT NULL REFERENCES square_merchants (square_merchant_id),
  square_payment_id TEXT PRIMARY KEY,
  square_order_id TEXT NOT NULL,
  square_customer_id TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE square_payments;

-- +goose StatementEnd
