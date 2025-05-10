-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_items (
  tenant_uuid UUID REFERENCES tenants (tenant_uuid),
  order_uuid UUID REFERENCES orders (order_uuid),
  order_item_uuid UUID PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE order_items;

-- +goose StatementEnd
