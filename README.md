# gamesync

gamesync is used to sync game save files to a server.

The project is seperated into client and server. 

The client is written in golang and syncs to the server using rsync.

Snapshots can be created on the remote via the client using restic. 

## Features

- Sync game save files with server
- Multiple devices
- Multiple users
- Snapshots

## Example usage

```sh
gamesync saves add openttd ~/.local/share/openttd/save
gamesync push openttd
gamesync snapshot create openttd
```

## Client environmental variables

- __GAMESYNC_CONFIG__: sets the config file to use, does not overwrite --config flag
- __GAMESYNC_STATE__: sets the state dir where sync states are stored

## Container environmental variables

- __GAMESYNC_LOOP__
    - when set to a number above 0 it will create users that often in seconds
    - if unset or below 1, only creates user on startup
    - useful if you want to add or remove users without restarting container

## Copyright

Copyright (C) 2025-2026 Patrik Bures
