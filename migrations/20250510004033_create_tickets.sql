-- +goose Up
-- +goose StatementBegin
CREATE TABLE tickets (
  tenant_uuid UUID REFERENCES tenants (tenant_uuid),
  order_item_uuid UUID REFERENCES order_items (order_item_uuid),
  ticket_uuid UUID PRIMARY KEY,
  checked_in BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE tickets;

-- +goose StatementEnd
