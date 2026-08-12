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
    repo_id BIGINT NOT NULL REFERENCES repos(repo_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_snapshots_parent ON snapshots(parent_snapshot_id);

-- +goose StatementBegin
CREATE FUNCTION delete_snapshot(p_snapshot_id BIGINT)
RETURNS BOOLEAN AS $$
DECLARE
    v_parent_id BIGINT;
BEGIN
    SELECT parent_snapshot_id INTO v_parent_id FROM snapshots
    WHERE snapshot_id = p_snapshot_id
    ;

    IF NOT FOUND THEN
        RETURN FALSE;
    END IF
    ;

    UPDATE snapshots
    SET parent_snapshot_id = v_parent_id
    WHERE parent_snapshot_id = p_snapshot_id
    ;

    -- only update branch head if parent is not null, which only
    -- occurs if it is deleting the first snapshot
    -- if it is the only snapshot in the branch, then the branch 
    -- will be deleted by the cascade
    IF v_parent_id IS NOT NULL THEN
        UPDATE branches
        SET head_snapshot_id = v_parent_id
        WHERE head_snapshot_id = p_snapshot_id
        ;
    END IF
    ;

    DELETE FROM snapshots
    WHERE snapshot_id = p_snapshot_id
    ;

    RETURN TRUE
    ;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd


CREATE TABLE branches
(
    branch_id BIGSERIAL NOT NULL PRIMARY KEY,
    repo_id BIGINT NOT NULL REFERENCES repos(repo_id) ON DELETE CASCADE,
    head_snapshot_id BIGINT NOT NULL REFERENCES snapshots(snapshot_id) ON DELETE CASCADE,
    branch_name VARCHAR(25) NOT NULL,

    UNIQUE (repo_id, branch_name)
);



CREATE TABLE files
(
    -- file hash is created by hashing all the raw chunk hashes in order, without any seperator.
    bytes BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
    bytes INT NOT NULL,
    bytes_compressed INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    chunk_hash BYTEA NOT NULL PRIMARY KEY,
    pending_deletion_at TIMESTAMPTZ NULL DEFAULT NULL
);

-- index only for chunks about to be deleted so listing it is fast
CREATE INDEX idx_chunks_pending_deletion
ON CHUNKS (pending_deletion_at, chunk_hash)
WHERE pending_deletion_at IS NOT NULL;

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
DROP INDEX idx_chunks_pending_deletion;

DROP TRIGGER trg_reparent_snapshots ON snapshots CASCADE;
DROP FUNCTION handle_referenced_row_delete() CASCADE;
