every api path will have a `api/v1` prefix

# user
they do not have passwords, just tokens are used

# some questions to answer:
- how to set default role
- should there be profiles? could each `sync_id` just associate with a `game_id` where you can have multiple `sync_id`s for each game?


| method | path                                                              | desc                                                       |
|--------|-------------------------------------------------------------------|------------------------------------------------------------|
| POST   | /user                                                             | create new user                                            |
| GET    | /user                                                             | get all users                                              |
| DELETE | /user/{user_id}                                                   | delete user                                                |
| GET    | /user/{user_id}                                                   | get specific user                                          |
| PUT    | /user/{user_id}/role/{role_id}                                    | change role of user                                        |
| GET    | /user/{user_id}/sync/{sync_id}                                    | get info about sync of game. storage used, chunk qty,      |
| POST   | /user/{user_id}/sync/{sync_id}/profile                            | create new profile for game                                |
| GET    | /user/{user_id}/sync/{sync_id}/profile/{profile_id}               | list snapshots for profile + some other info like its name |
| GET    | /user/{user_id}/sync/{sync_id}/profile/{profile_id}/{snapshot_id} | get all chunks associated with that snapshot               |
| POST   | /role                                                             | create new role                                            |
| GET    | /role                                                             | get all roles                                              |
| GET    | /role/{role_id}                                                   | get permissions for specific role                          |
| PUT    | /role/{role_id}                                                   | set permissions for role, defined in json data             |
| GET    | /game                                                             | get all games                                              |
| GET    | /game/{game_id}                                                   | get info of specific game                                  |
| POST   | /chunk                                                            | do you have these chunks? returns missing chunks           |
| GET    | /chunk/{chunk_sha256}                                             | download specific chunk                                    |
| POST   | /chunk/{chunk_sha256}                                             | upload specific chunk                                      |
