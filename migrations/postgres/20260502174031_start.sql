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
    UNIQUE (user_id, role_id),

    FOREIGN KEY (role_id) REFERENCES roles(role_id)
);

CREATE TABLE tokens
(
    token_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token_hash BYTEA NOT NULL,

    UNIQUE (token_hash),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- +goose Down
DELETE users;
DELETE roles;
DELETE tokens;
