-- +goose Up

CREATE TABLE roles
(
    role_id SERIAL NOT NULL,
    role_name VARCHAR(25) NOT NULL,

    UNIQUE (role_id),
    UNIQUE (role_name)
);

CREATE TABLE users
(
    user_id BIGSERIAL NOT NULL,
    user_name VARCHAR(25) NOT NULL,
    role_id INTEGER NOT NULL,

    UNIQUE (user_id),
    UNIQUE (user_name),
    UNIQUE (role_id),

    FOREIGN KEY (role_id) REFERENCES roles(role_id)
);

-- +goose Down
DELETE users;
DELETE roles
