-- name: GetTenant :one
SELECT * FROM tenants
WHERE tenant_uuid = $1 LIMIT 1;

-- name: ListTenants :many
SELECT * FROM tenants;

-- name: CreateTenant :one
INSERT INTO tenants (tenant_uuid, name)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateTenant :one
UPDATE tenants
SET name = $2
WHERE tenant_uuid = $1
RETURNING *;

-- name: DeleteTenant :one
DELETE FROM tenants
WHERE tenant_uuid = $1
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE user_uuid = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (user_uuid, username, name, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = $2, name = $3, password_hash = $4
WHERE user_uuid = $1
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE user_uuid = $1
RETURNING *;
