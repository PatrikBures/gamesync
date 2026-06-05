-- name: ListPermissions :many
SELECT * FROM permissions
;

-- name: ListPermsWithNames :many
SELECT * FROM permissions
WHERE perm_name = ANY(sqlc.arg(perm_names)::VARCHAR[])
;

-- name: InsertPermission :exec
INSERT INTO permissions (perm_id, perm_name)
VALUES ($1, $2)
;

-- name: UpdatePermissionName :exec
UPDATE permissions
SET perm_name = $2
WHERE perm_id = $1
;

-- name: DeletePermission :execrows
DELETE FROM permissions
WHERE perm_id = $1
;
