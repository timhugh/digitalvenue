-- +goose Up
-- +goose StatementBegin
CREATE TABLE tenants (
  tenant_uuid UUID PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW (),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE tenants;

-- +goose StatementEnd
