-- name: ListBranches :many
SELECT * FROM branches
WHERE repo_id = (
    SELECT repo_id FROM repos 
    WHERE repo_name = $1
)
;

-- name: GetBranchWithName :one
SELECT * FROM branches
WHERE repo_id = $1
AND branch_name = $2
LIMIT 1
;

-- name: CreateBranch :exec
INSERT INTO branches (repo_id, branch_name)
VALUES ($1, $2)
;

-- name: DeleteBranch :exec
DELETE FROM branches
WHERE repo_id = $1 AND branch_id = $2
;
