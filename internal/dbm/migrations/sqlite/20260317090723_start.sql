-- +goose Up

CREATE TABLE permission
(
    permission_id INTEGER NOT NULL,
    permission_name TEXT NOT NULL,

    PRIMARY KEY (permission_id),
    UNIQUE (permission_name)
);

CREATE TABLE role
(
    role_id INTEGER NOT NULL,
    role_name TEXT NOT NULL,

    PRIMARY KEY (role_id),
    UNIQUE (role_name)
);

CREATE TABLE role_permission
(
    role_permission_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,

    PRIMARY KEY (role_permission_id),
    UNIQUE (role_id, permission_id),
    FOREIGN KEY (permission_id) REFERENCES permission(permission_id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES role(role_id) ON DELETE CASCADE
);

CREATE TABLE user
(
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    user_name TEXT NOT NULL,

    PRIMARY KEY (user_id),
    UNIQUE (user_name),
    FOREIGN KEY (role_id) REFERENCES role(role_id)
);

CREATE TABLE ssh_key
(
    key_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    fingerprint TEXT NOT NULL,
    pk BLOB NOT NULL,
    comment TEXT NOT NULL,

    PRIMARY KEY (key_id),
    UNIQUE (fingerprint),
    FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE
);

CREATE TABLE game
(
    game_id INTEGER NOT NULL,
    canonical_name TEST NOT NULL,
    slug TEXT NOT NULL,

    PRIMARY KEY (game_id),
    UNIQUE (canonical_name),
    UNIQUE (slug)
);

CREATE TABLE user_game
(
    user_id INTEGER NOT NULL,
    game_id INTEGER NOT NULL,
    alias INTEGER NOT NULL,

    PRIMARY KEY (user_id, game_id),
    UNIQUE (alias),
    FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES game(game_id)
);
-- +goose Down
