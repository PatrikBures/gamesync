
-- name: CommitToBranch :one
WITH current_branch AS (
    SELECT head_snapshot_id
    FROM branches
    WHERE branches.repo_id = $1 AND branches.branch_name = $2
),
new_snapshot AS (
    INSERT INTO snapshots (repo_id, parent_snapshot_id)
    VALUES (
        $1,
        (SELECT head_snapshot_id FROM current_branch)
    )
    RETURNING *
),
upsert_branch AS (
    INSERT INTO branches (repo_id, branch_name, head_snapshot_id)
    SELECT $1, $2, new_snapshot.snapshot_id
    FROM new_snapshot
    ON CONFLICT (repo_id, branch_name)
    DO UPDATE SET head_snapshot_id = EXCLUDED.head_snapshot_id
    RETURNING *
)
SELECT 
    new_snapshot.snapshot_id,
    new_snapshot.parent_snapshot_id,
    upsert_branch.branch_id
FROM new_snapshot CROSS JOIN upsert_branch
;


-- name: HasAncestor :one
-- checks if 2 is ancestor of 1
SELECT snapshot_has_ancestor($1, $2)
;


-- name: ListFileHashes :many
SELECT file_hash FROM files
WHERE file_hash = ANY(sqlc.arg(file_hash)::BYTEA[])
;

-- name: CreateFile :exec
INSERT INTO files (file_hash, bytes)
VALUES (
    $1,
    (
        SELECT COALESCE(SUM(bytes), 0) FROM chunks
        WHERE chunk_hash = ANY(sqlc.arg(chunk_hash)::BYTEA[])
    )
)
;

-- name: ConnectFileWithChunks :copyfrom
-- connects files with chunks with the chunk order
INSERT INTO file_chunks (file_hash, chunk_hash, chunk_order)
VALUES ($1, $2, $3)
;

-- name: ConnectSnapshotWithFiles :copyfrom
INSERT INTO snapshot_files (file_hash, snapshot_id, file_path)
VALUES ($1, $2, $3)
;


-- name: UpdateBranchHead :exec
UPDATE branches
SET head_snapshot_id = $2
WHERE branch_id = $1
;

-- name: GetSnapshot :one
SELECT * FROM snapshots
WHERE snapshot_id = $1
AND repo_id = $2
;


-- name: GetSnapshotFiles :many
SELECT * FROM snapshot_files
WHERE snapshot_id = $1
;

-- name: GetFileChunkHashes :many
SELECT chunk_hash FROM file_chunks
WHERE file_hash = $1
ORDER BY chunk_order
;


-- name: DeleteSnapshot :one
SELECT delete_snapshot($1) AS was_deleted
;
