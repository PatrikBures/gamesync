-- name: CreateSnapshot :one
INSERT INTO snapshots (parent_snapshot_id)
VALUES (
    (
        SELECT head_snapshot_id FROM branches 
        WHERE repo_id = (
            SELECT repo_id FROM repos
            WHERE user_id = $1
            AND repo_name = $2
            LIMIT 1
        )
        AND branch_name = $3
        LIMIT 1
    )
)
RETURNING snapshot_id
;

