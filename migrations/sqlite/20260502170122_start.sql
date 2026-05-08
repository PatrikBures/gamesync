-- +goose Up

CREATE TABLE roles
(
    role_id INTEGER PRIMARY KEY,
    role_name TEXT NOT NULL,

    UNIQUE (role_name)
);

CREATE TABLE users
(
    user_id INTEGER PRIMARY KEY,
    user_name TEXT NOT NULL,
    role_id INTEGER NOT NULL,

    UNIQUE (user_name),
    UNIQUE (user_id, role_id),

    FOREIGN KEY (role_id) REFERENCES roles(role_id)
);

CREATE TABLE tokens
(
    user_id INTEGER PRIMARY KEY,
    token BLOB NOT NULL,

    UNIQUE (token),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- +goose Down
DELETE users;
DELETE roles;
