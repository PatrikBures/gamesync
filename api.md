every api path will have a `api/v1` prefix

# user
they do not have passwords, just tokens are used

# sync
a sync has the following:
- sync_id
- game_id

# game
a game has following info
- full name
- slug


# some questions to answer:
- how to set default role
  just have an env var for it


| Method | Path                                                    | Description                                                                                         |
|--------|---------------------------------------------------------|-----------------------------------------------------------------------------------------------------|
| POST   | /chunks                                                 | do you have these chunks? returns missing chunks                                                    |
| GET    | /chunks/{sha256}                                        | download specific chunk                                                                             |
| POST   | /chunks/{sha256}                                        | upload specific chunk                                                                               |
| POST   | /games                                                  | create new game                                                                                     |
| GET    | /games                                                  | get all games                                                                                       |
| DELETE | /games/{gid}                                            | delete game                                                                                         |
| GET    | /games/{gid}                                            | get info about specific game                                                                        |
| PUT    | /games/{gid}                                            | update game information like name                                                                   |
| PUT    | /games/{gid}/merge                                      | merge multiple gids into provided gid                                                               |
| POST   | /roles                                                  | create new role                                                                                     |
| GET    | /roles                                                  | list all roles                                                                                      |
| GET    | /roles/{role_id}                                        | get permissions for role                                                                            |
| PUT    | /roles/{role_id}                                        | set permissions for role, defined in json data                                                      |
| POST   | /users                                                  | create new user, returning their token                                                              |
| GET    | /users                                                  | list all users                                                                                      |
| DELETE | /users/{uid}                                            | delete user                                                                                         |
| GET    | /users/{uid}                                            | get user info, including quota limit                                                                |
| PUT    | /users/{uid}/name                                       | update user name                                                                                    |
| PUT    | /users/{uid}/role                                       | change role of user (role provided in body)                                                         |
| POST   | /users/{uid}/syncs                                      | create new sync, with name and optional gid                                                         |
| GET    | /users/{uid}/syncs                                      | list all syncs (optionally filter using query?)                                                     |
| DELETE | /users/{uid}/syncs/{sid}                                | delete sync                                                                                         |
| GET    | /users/{uid}/syncs/{sid}                                | get info about sync of game. storage used, chunk qty, profile list, and latest snapstot per profile |
| PUT    | /users/{uid}/syncs/{sid}                                | update sync name or game                                                                            |
| POST   | /users/{uid}/syncs/{sid}/profiles                       | create profile providing its name. default "main" exists on sync creation                           |
| DELETE | /users/{uid}/syncs/{sid}/profiles                       | delete profile                                                                                      |
| GET    | /users/{uid}/syncs/{sid}/profiles                       | list profiles of a sync                                                                             |
| POST   | /users/{uid}/syncs/{sid}/profiles/{pid}/snapshots       | create new snapshot, providing a list of chunks with their path, order and sha256                   |
| GET    | /users/{uid}/syncs/{sid}/profiles/{pid}/snapshots       | list snapshots                                                                                      |
| DELETE | /users/{uid}/syncs/{sid}/profiles/{pid}/snapshots/{sid} | delete snapshot (will not delete chunks?)                                                           |
| GET    | /users/{uid}/syncs/{sid}/profiles/{pid}/snapshots/{sid} | list chunks associated with snapshot + metadata                                                     |



