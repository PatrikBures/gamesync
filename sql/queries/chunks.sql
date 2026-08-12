-- name: CreateChunk :exec
INSERT INTO chunks (chunk_hash, bytes, bytes_compressed)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING
;

-- name: CheckChunk :one
SELECT COUNT(*) > 0 FROM chunks
WHERE chunk_hash = $1
LIMIT 1
;

-- name: GetChunkHashes :many
SELECT chunk_hash FROM chunks
WHERE chunk_hash = ANY(sqlc.arg(chunk_hash)::BYTEA[])
ORDER BY chunk_hash ASC
;


-- name: GetChunkHashesClearMark :many
WITH existing_chunks AS (
    SELECT chunk_hash FROM chunks
    WHERE chunk_hash = ANY(sqlc.arg(chunk_hash)::BYTEA[])
    ORDER BY chunk_hash ASC
), 
updated AS (
    UPDATE chunks
    SET pending_deletion_at = NULL
    FROM existing_chunks
    WHERE chunks.chunk_hash = existing_chunks.chunk_hash
    AND chunks.pending_deletion_at IS NOT NULL
)
SELECT chunk_hash FROM existing_chunks
;
