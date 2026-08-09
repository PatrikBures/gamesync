-- +goose Up
CREATE TABLE repos
(
    repo_id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    repo_name VARCHAR(25) NOT NULL,

    UNIQUE (user_id, repo_name)
);

CREATE TABLE snapshots
(
    snapshot_id BIGSERIAL NOT NULL PRIMARY KEY,
    parent_snapshot_id BIGINT REFERENCES snapshots(snapshot_id),
    repo_id BIGINT NOT NULL REFERENCES repos(repo_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_snapshots_parent ON snapshots(parent_snapshot_id);

CREATE TABLE branches
(
    branch_id BIGSERIAL NOT NULL PRIMARY KEY,
    repo_id BIGINT NOT NULL REFERENCES repos(repo_id) ON DELETE CASCADE,
    branch_name VARCHAR(25) NOT NULL,
    head_snapshot_id BIGINT NULL REFERENCES snapshots(snapshot_id),

    UNIQUE (repo_id, branch_name)
);



CREATE TABLE files
(
    -- file hash is created by hashing all the raw chunk hashes in order, without any seperator. 
    bytes BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    file_hash BYTEA NOT NULL PRIMARY KEY
);

CREATE TABLE snapshot_files
(
    snapshot_id BIGINT NOT NULL REFERENCES snapshots(snapshot_id) ON DELETE CASCADE,
    file_hash BYTEA NOT NULL REFERENCES files(file_hash) ON DELETE RESTRICT,
    file_path VARCHAR(500) NOT NULL,

    PRIMARY KEY (snapshot_id, file_path)
);
CREATE INDEX idx_snapshot_files_file_hash ON snapshot_files(file_hash);

CREATE TABLE chunks
(
    bytes BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    chunk_hash BYTEA NOT NULL PRIMARY KEY
);

CREATE TABLE file_chunks
(
    chunk_order INT NOT NULL,
    file_hash BYTEA NOT NULL REFERENCES files(file_hash) ON DELETE CASCADE,
    chunk_hash BYTEA NOT NULL REFERENCES chunks(chunk_hash),

    PRIMARY KEY (file_hash, chunk_order)
);
CREATE INDEX idx_file_chunks_chunk_hash ON file_chunks(chunk_hash);



-- +goose Down
DELETE repos;
DELETE snapshots;
DELETE branches;
DELETE files;
DELETE snapshot_files;
DELETE chunks;
DELETE file_chunks;

DROP INDEX idx_snapshots_parent;
DROP INDEX idx_snapshot_files_file_hash;
DROP INDEX idx_file_chunks_chunk_hash;
