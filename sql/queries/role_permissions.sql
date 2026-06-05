-- name: DeleteRolePermsLT :execrows
DELETE FROM role_permissions
WHERE role_id < $1
;

-- name: InsertRolePerms :copyfrom
INSERT INTO role_permissions (role_id, perm_id)
VALUES ($1, $2)
;

-- name: ListRolePerms :many
SELECT * FROM role_permissions
WHERE role_id = $1
;

-- name: ListRolePermNamesWithName :many
-- returns slice of all the permissions a role has
--
-- the returned permissions are their names
SELECT perm_name FROM permissions
JOIN role_permissions USING (perm_id) 
WHERE role_permissions.role_id = $1
;
