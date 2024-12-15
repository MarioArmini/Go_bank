-- name: CreateTransfer :one
INSERT INTO "Transfers"(
    "senderId",
    "recipientId",
    amount
) VALUES(
    $1, $2, $3
) RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM "Transfers" WHERE "Id" = $1;

-- name: GetTransfer :one
SELECT * FROM "Transfers"
WHERE "Id" = $1
LIMIT 1;

-- name: GetAllTransfers :many
SELECT * FROM "Transfers"
ORDER BY "Id"
LIMIT $1
OFFSET $2;