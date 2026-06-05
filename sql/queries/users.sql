-- name: ListUsers :many
SELECT * FROM users
ORDER BY user_name;

-- name: InsertUser :one
INSERT INTO users (user_name, role_id)
VALUES ($1, $2)
RETURNING user_id
;

-- name: InsertUserWithId :exec
INSERT INTO users (user_id, user_name, role_id)
VALUES ($1, $2, $3)
;

-- name: GetUserWithName :one
SELECT * FROM users
WHERE user_name = $1
;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1
;

-- name: UpdateUserName :exec
UPDATE users
SET user_name = $2
WHERE user_id = $1
;
