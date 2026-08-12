-- this is where all garbage collection related quereies are put


-- name: GarbageDeleteOrphanedFiles :many
WITH orphaned_files AS (
    SELECT files.file_hash FROM files
    WHERE files.created_at < NOW() - (sqlc.arg(older_than_microseconds)::BIGINT * INTERVAL '1 microseconds')
    AND NOT EXISTS(
        SELECT 1
        FROM snapshot_files
        WHERE snapshot_files.file_hash = files.file_hash
    )
    LIMIT 500
    FOR UPDATE OF files SKIP LOCKED
)
DELETE FROM files
WHERE file_hash IN (SELECT file_hash FROM orphaned_files)
RETURNING file_hash
;

-- name: GarbageMarkOrphanedChunks :one
WITH to_mark AS (
    SELECT chunks.chunk_hash FROM chunks
    WHERE chunks.pending_deletion_at IS NULL
    AND chunks.created_at < NOW() - (sqlc.arg(older_than_microseconds)::BIGINT * INTERVAL '1 microseconds')
    AND NOT EXISTS (
        SELECT 1 FROM file_chunks
        WHERE file_chunks.chunk_hash = chunks.chunk_hash
    )
    LIMIT 1000
    FOR UPDATE OF chunks SKIP LOCKED
),
updated AS (
    UPDATE chunks
    SET pending_deletion_at = NOW()
    WHERE chunk_hash IN (SELECT chunk_hash FROM to_mark)
    RETURNING 1
)
SELECT COUNT(*) AS marked_count FROM updated;
;


-- name: GarbageListMarkedChunks :many
SELECT chunks.chunk_hash FROM chunks
WHERE chunks.pending_deletion_at < NOW() - (sqlc.arg(older_than_microseconds)::BIGINT * INTERVAL '1 microseconds')
AND NOT EXISTS (
    SELECT 1 FROM file_chunks
    WHERE file_chunks.chunk_hash = chunks.chunk_hash
)
ORDER BY chunks.pending_deletion_at ASC, chunks.chunk_hash ASC
LIMIT 100
FOR UPDATE OF chunks SKIP LOCKED
;


-- name: GarbageDeleteChunks :execrows
DELETE FROM chunks
WHERE chunk_hash = ANY(sqlc.arg(chunk_hash)::BYTEA[])
;
