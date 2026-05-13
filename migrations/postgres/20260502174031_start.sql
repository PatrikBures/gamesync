-- +goose Up

CREATE TABLE roles
(
    role_id SERIAL NOT NULL,
    role_name VARCHAR(25) NOT NULL,

    PRIMARY KEY (role_id),
    UNIQUE (role_id),
    UNIQUE (role_name)
);

CREATE TABLE users
(
    user_id BIGSERIAL NOT NULL,
    user_name VARCHAR(25) NOT NULL,
    role_id INTEGER NOT NULL,

    PRIMARY KEY (user_id),
    UNIQUE (user_id),
    UNIQUE (user_name),
    FOREIGN KEY (role_id) REFERENCES roles(role_id)
);

CREATE TABLE tokens
(
    token_id BIGSERIAL,
    user_id BIGINT NOT NULL,
    token_hash BYTEA NOT NULL,

    PRIMARY KEY (token_id),
    UNIQUE (token_hash),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE permissions
(
    perm_id SERIAL,
    perm_name VARCHAR(50) NOT NULL,

    PRIMARY KEY (perm_id),
    UNIQUE (perm_name)
);

CREATE TABLE role_permissions
(
    role_perm_id SERIAL,
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
