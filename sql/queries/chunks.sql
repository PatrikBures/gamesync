-- name: CreateChunk :exec
INSERT INTO chunks (chunk_hash, bytes)
VALUES ($1, $2)
;

-- name: CheckChunk :one
SELECT COUNT(*) > 0 FROM chunks
WHERE chunk_hash = $1
LIMIT 1
;
