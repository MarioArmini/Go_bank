-- name: CreateAccount :one
INSERT INTO "Accounts" (
    owner,
    balance,
    currency
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM "Accounts"
WHERE "Id" = $1 LIMIT 1;

-- name: GetAllAccounts :many
SELECT * FROM "Accounts"
ORDER BY "Id"
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE "Accounts"
SET balance = $2
WHERE "Id" = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM "Accounts" WHERE "Id" = $1;
