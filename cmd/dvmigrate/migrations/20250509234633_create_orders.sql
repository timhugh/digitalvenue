-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
  tenant_uuid UUID NOT NULL REFERENCES tenants (tenant_uuid),
  customer_uuid UUID,
  order_uuid UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
  'down SQL query';

-- +goose StatementEnd
