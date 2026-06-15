-- name: ListBranches :many
SELECT * FROM branches
WHERE repo_id = (
    SELECT repo_id FROM repos 
    WHERE repo_name = $1
)
;

-- name: CreateBranch :exec
INSERT INTO branches (repo_id, branch_name)
VALUES (
    (
        SELECT repo_id FROM repos
        WHERE repo_name = $1
    ),
    $2
)
;
