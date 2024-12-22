-- name: CreateAccount :one
INSERT INTO "Accounts" (
    owner,
    balance,
    currency,
    "interestRate"
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM "Accounts"
WHERE "Id" = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM "Accounts"
WHERE "Id" = $1 LIMIT 1
FOR NO KEY UPDATE;

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

-- name: GetBalance :one
SELECT balance FROM "Accounts"
WHERE "Id" = $1;
