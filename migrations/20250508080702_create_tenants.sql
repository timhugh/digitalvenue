-- +goose Up
-- +goose StatementBegin
CREATE TABLE tenants (tenant_uuid UUID PRIMARY KEY, name TEXT NOT NULL);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE tenants;

-- +goose StatementEnd
