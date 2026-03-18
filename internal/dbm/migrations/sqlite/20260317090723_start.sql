-- +goose Up
CREATE TABLE role
(
    role_id INTEGER NOT NULL,
    role_name TEXT NOT NULL,

    perm_game_delete        BOOLEAN CHECK(perm_game_delete          IN(0, 1)) DEFAULT 0,
    perm_game_delete_own    BOOLEAN CHECK(perm_game_delete_own      IN(0, 1)) DEFAULT 0,
    perm_game_rename        BOOLEAN CHECK(perm_game_rename          IN(0, 1)) DEFAULT 0,
    perm_game_rename_own    BOOLEAN CHECK(perm_game_rename_own      IN(0, 1)) DEFAULT 0,
    perm_game_add           BOOLEAN CHECK(perm_game_add             IN(0, 1)) DEFAULT 0,
    perm_user_add           BOOLEAN CHECK(perm_user_add             IN(0, 1)) DEFAULT 0,
    perm_user_delete        BOOLEAN CHECK(perm_user_delete          IN(0, 1)) DEFAULT 0,
    perm_user_rename        BOOLEAN CHECK(perm_user_rename          IN(0, 1)) DEFAULT 0,
    perm_user_rename_own    BOOLEAN CHECK(perm_user_rename_own      IN(0, 1)) DEFAULT 0,
    perm_user_change_role   BOOLEAN CHECK(perm_user_change_role     IN(0, 1)) DEFAULT 0,
    perm_user_list          BOOLEAN CHECK(perm_user_list            IN(0, 1)) DEFAULT 0,
    perm_role_create        BOOLEAN CHECK(perm_role_create          IN(0, 1)) DEFAULT 0,
    perm_role_remove        BOOLEAN CHECK(perm_role_remove          IN(0, 1)) DEFAULT 0,
    perm_role_change_perms  BOOLEAN CHECK(perm_role_change_perms    IN(0, 1)) DEFAULT 0,
    perm_sync               BOOLEAN CHECK(perm_sync                 IN(0, 1)) DEFAULT 0,

    PRIMARY KEY (role_id),
    UNIQUE (role_name)
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
