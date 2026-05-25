-- +goose Up
CREATE TABLE syncs
(
    sync_id INTEGER,
    user_id INTEGER NOT NULL,
    sync_name TEXT NOT NULL,

    PRIMARY KEY (sync_id),
    UNIQUE (user_id, sync_name),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);


CREATE TABLE sync_profiles
(
    sync_profile_id INTEGER,
    sync_id INTEGER NOT NULL,
    sync_profile_name TEXT NOT NULL,

    PRIMARY KEY (sync_profile_id),
    UNIQUE (sync_id, sync_profile_name),
    FOREIGN KEY (sync_id) REFERENCES syncs(sync_id)
);


-- +goose Down
DELETE syncs;
DELETE sync_profiles;
