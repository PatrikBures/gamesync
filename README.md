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
gamesync save add openttd ~/.local/share/openttd/save
gamesync push openttd
gamesync snapshot create openttd
```

## Set up wrap command for automatic syncs when starting and exiting games

### How the wrap command works

The command is formatted like this:
```sh
gamesync wrap GAME_ID -- COMMAND
```

Whatever is after "--" will be blindly ran as a command.

All that gamesync will do is pull before running the command, and push after the command exited.
And optionally create a snapshot with restic with "-s" or "-S" flags.

### Steam

1. Add the savepath to gamesync
2. Select the game you want to sync from your library
3. Right click and pick properties
4. In the general tab add this to your Launch Options
```sh
gamesync wrap GAME_ID -- %command%
```

### Lutris

1. Add the savepath to gamesync
2. Select the game you want to sync from your library
3. Right click and pick configure
4. Open the system options and scroll down to command prefix
5. Add this to your command prefix
```sh
gamesync wrap GAME_ID --
```

### .desktop file

1. Find the desktop file you want to modify. 
You can probably find your programs .desktop file located in /usr on Linux.
This command can help you find it:
```sh
find /usr -iname '*NAME*.desktop'
```
Replace NAME with the name of the app or game.

2. Then copy that file to your users application directory. Usually located at `~/.local/share/applications`

If your desktop file already is somewhere in your home dir, you should still create a copy of it with a suffix. 
Just in case you regret your modifications.

2. Open the copied .desktop file and modify the Exec command to something like this
```ini
Exec=gamesync wrap GAME_ID -- [whatever was originally here]
```

### Useful wrap flags

- __-n__: Sends notification whenever it pushed, pulled or created a snapshot. Notifying you if it succeded or failed.
- __-S__ Creates snapshot only if there are changes after pushing. 
- __--no-pull__: Useful if you want to just try out the wrap command or use as a backup without modifying your local save files.
- __-e__: Exits on error. For example, if the pull failed it exits, preventing the game from launching.

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
