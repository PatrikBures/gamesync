-- name: InsertToken :one
INSERT INTO tokens (user_id, token_hash)
VALUES ($1, $2)
RETURNING token_id
;

-- name: GetUserFromToken :one
SELECT * FROM users
WHERE user_id = (
    SELECT user_id FROM tokens
    WHERE token_hash = $1
)
;
