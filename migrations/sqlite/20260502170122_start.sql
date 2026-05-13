-- +goose Up

CREATE TABLE roles
(
    role_id INTEGER,
    role_name TEXT NOT NULL,

    PRIMARY KEY (role_id),
    UNIQUE (role_name)
);

CREATE TABLE users
(
    user_id INTEGER,
    user_name TEXT NOT NULL,
    role_id INTEGER NOT NULL,

    PRIMARY KEY (user_id),
    UNIQUE (user_id),
    UNIQUE (user_name),
    FOREIGN KEY (role_id) REFERENCES roles(role_id)
);

CREATE TABLE tokens
(
    token_id INTEGER,
    user_id INTEGER NOT NULL,
    token_hash BLOB NOT NULL,

    PRIMARY KEY (token_id),
    UNIQUE (token_hash),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE permissions
(
    perm_id INTEGER,
    perm_name TEXT NOT NULL,

    PRIMARY KEY (perm_id),
    UNIQUE (perm_name)
);

CREATE TABLE role_permissions
(
    role_perm_id INTEGER,
    role_id INTEGER NOT NULL,
    perm_id INTEGER NOT NULL,

    PRIMARY KEY (role_perm_id),
    UNIQUE (role_id, perm_id),
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    FOREIGN KEY (perm_id) REFERENCES permissions(perm_id) ON DELETE CASCADE
);

-- +goose Down
DELETE users;
DELETE roles;
DELETE tokens;
DELETE permissions;
DELETE role_permissions;
