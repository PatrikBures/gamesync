
-- name: CreateSnapshot :one
-- returns the new snapshot_id
--
-- accepts repo_id and branch_id to find the parent snapshot
-- using head snapshot id from branches
INSERT INTO snapshots (parent_snapshot_id)
VALUES (
    (
        SELECT head_snapshot_id FROM branches 
        WHERE repo_id = $1
        AND branch_id = $2
        LIMIT 1
    )
)
RETURNING snapshot_id
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
        SELECT SUM(bytes) FROM chunks
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
