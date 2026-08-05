-- name: ListRepos :many
SELECT * FROM repos
WHERE user_id = $1
;

-- name: GetRepoWithName :one
SELECT * FROM repos
WHERE user_id = $1
AND repo_name = $2
LIMIT 1
;

-- name: CreateRepo :one
INSERT INTO repos (user_id, repo_name)
VALUES ($1, $2)
RETURNING repo_id
;

-- name: DeleteRepo :exec
DELETE FROM repos
WHERE repo_id = $1
;
