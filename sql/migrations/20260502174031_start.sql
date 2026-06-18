-- +goose Up

CREATE TABLE roles
(
    role_id SERIAL NOT NULL PRIMARY KEY,
    role_name VARCHAR(25) NOT NULL UNIQUE
);

ALTER SEQUENCE roles_role_id_seq RESTART 100;


CREATE TABLE users
(
    user_id BIGSERIAL NOT NULL PRIMARY KEY,
    user_name VARCHAR(25) NOT NULL UNIQUE,
    role_id INTEGER NOT NULL REFERENCES roles(role_id)
);

CREATE TABLE tokens
(
    token_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    token_hash BYTEA NOT NULL UNIQUE
);

CREATE TABLE permissions
(
    perm_id SERIAL PRIMARY KEY,
    perm_name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE role_permissions
(
    role_perm_id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
    perm_id INTEGER NOT NULL REFERENCES permissions(perm_id) ON DELETE CASCADE,

    UNIQUE (role_id, perm_id)
);

-- +goose Down
DELETE users;
DELETE roles;
DELETE tokens;
DELETE permissions;
DELETE role_permissions;
