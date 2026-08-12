-- name: ListBranches :many
SELECT
    branches.branch_id,
    branches.repo_id,
    branches.head_snapshot_id,
    snapshots.parent_snapshot_id,
    branches.branch_name
FROM branches
INNER JOIN snapshots ON branches.head_snapshot_id = snapshots.snapshot_id
WHERE branches.repo_id = $1
AND (sqlc.narg('branchName')::text IS NULL OR branch_name = sqlc.narg('branchName'))
ORDER BY branch_id
;

-- name: GetBranchWithName :one
SELECT * FROM branches
WHERE repo_id = $1
AND branch_name = $2
LIMIT 1
;

-- name: CreateBranch :one
INSERT INTO branches (repo_id, branch_name, head_snapshot_id)
VALUES ($1, $2, $3)
RETURNING branch_id
;

-- name: DeleteBranch :exec
DELETE FROM branches
WHERE repo_id = $1 AND branch_id = $2
;
