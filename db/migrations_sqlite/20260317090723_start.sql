-- +goose Up
CREATE TABLE user_type
(
    user_type_id INTEGER NOT NULL,
    user_type_name TEXT NOT NULL,

    PRIMARY KEY (user_type_id),
    UNIQUE (user_type_name)
);

INSERT INTO user_type (user_type_name) VALUES ('admin'), ('user');

CREATE TABLE user
(
    user_id INTEGER NOT NULL,
    user_type_id INTEGER NOT NULL,
    user_name TEXT NOT NULL,

    PRIMARY KEY (user_id),
    UNIQUE (user_name),
    FOREIGN KEY (user_type_id) REFERENCES user_type(user_type_id)
);

CREATE TABLE ssh_key
(
    key_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    hash TEXT NOT NULL,
    data TEXT NOT NULL,

    PRIMARY KEY (key_id),
    UNIQUE (hash),
    FOREIGN KEY (user_id) REFERENCES user(user_id)
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
    FOREIGN KEY (user_id) REFERENCES user(user_id),
    FOREIGN KEY (game_id) REFERENCES game(game_id)
);
-- +goose Down
