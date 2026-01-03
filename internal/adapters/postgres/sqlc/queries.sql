-- name: ListProducts :many
SELECT * FROM products;

-- name: FindProductByID :one
SELECT * FROM products WHERE id = $1;


-- name: CreateOrder :one
INSERT INTO orders (
    customer_id
) VALUES ($1) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_cents, status)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: FindProductByIDForUpdate :one
SELECT * FROM products WHERE id = $1 FOR UPDATE;

-- name: UpdateOrderItemStatus :one
UPDATE order_items
SET status = $2
WHERE id = $1
RETURNING *;

-- name: UpdateProductQuantity :one
UPDATE products
SET quantity = $2
WHERE id = $1
RETURNING *;
