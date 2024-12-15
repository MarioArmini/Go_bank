-- name: CreateEntry :one
INSERT INTO "Entries"(
    "accountId",
    amount
) VALUES(
    $1, $2
) RETURNING *;

-- name: UpdateEntry :one
UPDATE "Entries"
SET amount = $2
WHERE "Id" = $1
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM "Entries" WHERE "Id" = $1;

-- name: GetEntry :one
SELECT * FROM "Entries" 
WHERE "Id" = $1
LIMIT 1;

-- name: GetAllEntries :many
SELECT * FROM "Entries"
ORDER BY "Id"
LIMIT $1
OFFSET $2;
