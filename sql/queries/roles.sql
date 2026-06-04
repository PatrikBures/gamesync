-- name: ListRoles :many
SELECT * FROM roles
;

-- name: InsertRole :one
INSERT INTO roles (role_name)
VALUES ($1)
RETURNING role_id
;

-- name: UpdateRoleName :exec
UPDATE roles
SET role_name = $2
WHERE role_id = $1
;

-- name: InsertRoleWithId :exec
INSERT INTO roles (role_id, role_name)
VALUES ($1, $2)
;

-- name: GetRoleWithId :one
SELECT * FROM roles
WHERE role_id = $1
;

-- name: GetRoleWithIdCount :one
SELECT COUNT(*) FROM roles
WHERE role_id = $1
;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE role_id = $1
;
