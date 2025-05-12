-- name: GetTenant :one
SELECT * FROM tenants
WHERE tenant_uuid = $1 LIMIT 1;

-- name: ListTenants :many
SELECT * FROM tenants;

-- name: CreateTenant :one
INSERT INTO tenants (tenant_uuid, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE user_uuid = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (user_uuid, username, name, password_hash, password_salt)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSquareMerchant :one
SELECT * FROM square_merchants
WHERE square_merchant_id = $1 LIMIT 1;

-- name: CreateSquareMerchant :one
INSERT INTO square_merchants (tenant_uuid, square_merchant_id, square_api_token, square_webhook_signature_key, ticketable_categories, location_ids)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSquarePayment :one
SELECT * FROM square_payments
WHERE square_payment_id = $1 LIMIT 1;

-- name: CreateSquarePayment :one
INSERT INTO square_payments (square_merchant_id, square_payment_id, square_order_id, square_customer_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCustomer :one
SELECT * FROM customers
WHERE customer_uuid = $1 LIMIT 1;

-- name: CreateCustomer :one
INSERT INTO customers (tenant_uuid, customer_uuid, given_name, family_name, email_address, phone_number)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListCustomers :many
SELECT * FROM customers
WHERE tenant_uuid = $1;

-- name: GetSquareCustomer :one
SELECT * FROM square_customers
WHERE square_customer_id = $1 LIMIT 1;

-- name: CreateSquareCustomer :one
INSERT INTO square_customers (square_merchant_id, square_customer_id, given_name, family_name, email_address, phone_number)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateSquareCustomer :one
UPDATE square_customers
SET given_name = $2, family_name = $3, email_address = $4, phone_number = $5, customer_uuid = $6
WHERE square_customer_id = $1
RETURNING *;

-- name: CreateOrder :one
INSERT INTO orders (tenant_uuid, customer_uuid, order_uuid)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (tenant_uuid, order_uuid, order_item_uuid, name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateTicket :one
INSERT INTO tickets (tenant_uuid, order_item_uuid, ticket_uuid, qr_code_url)
VALUES ($1, $2, $3, $4)
RETURNING *;
