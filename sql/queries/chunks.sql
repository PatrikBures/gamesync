-- name: CreateChunk :exec
INSERT INTO chunks (chunk_hash, bytes)
VALUES ($1, $2)
;
