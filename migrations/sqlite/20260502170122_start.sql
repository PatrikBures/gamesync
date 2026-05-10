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
    token_id INTEGER PRIMARY KEY,
    user_id NOT NULL,
    token_hash BLOB NOT NULL,

    UNIQUE (token_hash),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE permissions
(
    perm_id INTEGER PRIMARY KEY,
    perm_name TEXT NOT NULL,

    UNIQUE (perm_name)
);

CREATE TABLE role_permissions
(
    role_perm_id INTEGER PRIMARY KEY,
    role_id INTEGER NOT NULL,
    perm_id INTEGER NOT NULL,

    UNIQUE (role_id, perm_id),
    FOREIGN KEY (role_id) REFERENCES roles(role_id),
    FOREIGN KEY (perm_id) REFERENCES roles(perm_id)
);


-- +goose Down
DELETE users;
DELETE roles;
DELETE tokens;
DELETE permissions;
DELETE role_permissions;
