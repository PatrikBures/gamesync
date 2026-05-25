-- +goose Up
CREATE TABLE syncs
(
    sync_id BIGSERIAL NOT NULL,
    user_id BIGINT NOT NULL,
    sync_name VARCHAR(25) NOT NULL,

    PRIMARY KEY (sync_id),
    UNIQUE (user_id, sync_name),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE snapshots
(
    snapshot_id BIGSERIAL NOT NULL,
    parent_snapshot_id BIGINT NULL,

    PRIMARY KEY (snapshot_id),
    FOREIGN KEY (parent_snapshot_id) REFERENCES snapshots(snapshot_id)
);

CREATE TABLE sync_profiles
(
    sync_profile_id BIGSERIAL NOT NULL,
    sync_id BIGINT NOT NULL,
    profile_name VARCHAR(25) NOT NULL,
    head_snapshot_id BIGSERIAL NOT NULL REFERENCES snapshots(snapshot_id),

    PRIMARY KEY (sync_profile_id),
    UNIQUE (sync_id, profile_name),
    FOREIGN KEY (sync_id) REFERENCES syncs(sync_id)
);



CREATE TABLE files
(
    file_hash BYTEA NOT NULL PRIMARY KEY,
    bytes BIGINT NOT NULL
);

CREATE TABLE snapshot_files
(
    file_hash BYTEA NOT NULL REFERENCES files(file_hash),
    snapshot_id BIGINT NOT NULL REFERENCES snapshots(snapshot_id),
    file_path VARCHAR(500) NOT NULL,
    PRIMARY KEY (file_hash, snapshot_id)
);


CREATE TABLE chunks
(
    chunk_hash BYTEA NOT NULL,
    PRIMARY KEY (chunk_hash)
);

CREATE TABLE file_chunks
(
    file_hash BYTEA NOT NULL REFERENCES files(file_hash),
    chunk_hash BYTEA NOT NULL REFERENCES chunks(chunk_hash),
    chunk_order SMALLINT NOT NULL,
    PRIMARY KEY (file_hash, chunk_order)
);




-- +goose Down
DELETE syncs;
DELETE sync_profiles;

