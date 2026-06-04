-- name: GetUserIdFromToken :one
SELECT user_id FROM tokens
WHERE token_hash = $1
;


-- name: InsertToken :one
INSERT INTO tokens (user_id, token_hash)
VALUES ($1, $2)
RETURNING token_id
;
