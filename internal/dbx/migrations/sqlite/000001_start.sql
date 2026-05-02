CREATE TABLE users
(
    user_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    role_id INTEGER NOT NULL,

    UNIQUE (name),
    UNIQUE (role_id),

    FOREIGN KEY (role_id) REFERENCES roles(role_id)
);

CREATE TABLE roles
(
    role_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,

    UNIQUE (name)
);
