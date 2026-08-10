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
