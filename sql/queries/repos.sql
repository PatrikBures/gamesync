-- name: ListRepos :many
SELECT * FROM repos
WHERE user_id = $1
;

-- name: CreateRepo :exec
INSERT INTO repos (user_id, repo_name)
VALUES ($1, $2)
;
