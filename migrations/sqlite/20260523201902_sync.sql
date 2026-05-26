-- +goose Up
CREATE TABLE repos
(
    repo_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    repo_name TEXT NOT NULL,

    PRIMARY KEY (repo_id),
    UNIQUE (user_id, repo_name),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE snapshots
(
    snapshot_id INTEGER NOT NULL,
    parent_snapshot_id INTEGER NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),

    PRIMARY KEY (snapshot_id),
    FOREIGN KEY (parent_snapshot_id) REFERENCES snapshots(snapshot_id)
);
CREATE INDEX idx_snapshots_parent ON snapshots(parent_snapshot_id);

CREATE TABLE branches
(
    branch_id INTEGER NOT NULL,
    repo_id INTEGER NOT NULL,
    branch_name TEXT NOT NULL,
    head_snapshot_id INTEGER NULL REFERENCES snapshots(snapshot_id),

    PRIMARY KEY (branch_id),
    UNIQUE (repo_id, branch_name),
    FOREIGN KEY (repo_id) REFERENCES repos(repo_id)
);



CREATE TABLE files
(
    file_hash BLOB NOT NULL PRIMARY KEY,
    bytes INTEGER NOT NULL
);

CREATE TABLE snapshot_files
(
    file_hash BLOB NOT NULL REFERENCES files(file_hash),
    snapshot_id INTEGER NOT NULL REFERENCES snapshots(snapshot_id),
    file_path TEXT NOT NULL,
    PRIMARY KEY (snapshot_id, file_path)
);
CREATE INDEX idx_snapshot_files_file_hash ON snapshot_files(file_hash);

CREATE TABLE chunks
(
    chunk_hash BLOB NOT NULL PRIMARY KEY
);

CREATE TABLE file_chunks
(
    file_hash BLOB NOT NULL REFERENCES files(file_hash),
    chunk_hash BLOB NOT NULL REFERENCES chunks(chunk_hash),
    chunk_order INTEGER NOT NULL,
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

