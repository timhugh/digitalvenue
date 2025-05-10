-- +goose Up
-- +goose StatementBegin
CREATE TABLE customers (
  tenant_uuid UUID REFERENCES tenants (tenant_uuid),
  customer_uuid UUID PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE customers;

-- +goose StatementEnd
