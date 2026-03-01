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

## Set up config

the config file is read from
```sh
$XDG_CONFIG_HOME/gamesync/config.yml
```
or
```sh
$HOME/.config/gamesync/config.yml
```
based on if __XDG_CONFIG_HOME__ is set

For macs it should be located at
```sh
$HOME/Library/Application/gamesync/config.yml
```

when gamesync is ran it will create the dir but not the config.yml which needs to be made manually

this is the required content of config.yml
```yml
server:
  host: host or ip of the server
  user: your user name on the server
  port: the ssh port
  identity_file: full path to your ssh key
```

i would suggest to not write any comments or something as a command will overwrite it. 

## Add games to config

### Manual way

You can manually type add the games to the config

example content of config.yml games:

```yml
games:
- id: game_id_1
  save_path: /full path/to/game 1
- id: game_id_2
  save_path: /full path/to/game 2
```

the game id is the id you will use in other commands to sync your save files, like the wrap command explained later

### Using save add command

```sh
gamesync save add game_id_1 "/full path/to/game 1"
gamesync save add game_id_2 "/full path/to/game 2"
```

You can use ~ to refer to your home directory like so:
```sh
gamesync save add game_id_1 ~/.local/share/game_1/save
```
_Your shell should convert that into a full path, if not use relative paths or full paths._

And you can also use relative paths
```sh
cd ~/.local/share
gamesync save add game_id_1 game_1/save
```
This will be converted into a full path by gamesync which will be added to config.yml


## Set up wrap command for automatic syncs when starting and exiting games

### How the wrap command works

The command is formatted like this:
```sh
gamesync wrap GAME_ID -- COMMAND
```

Whatever is after "--" will be blindly ran as a command.

All that gamesync will do is pull before running the command, and push after the command exited.
And optionally create a snapshot using restic with the -s or -S flags.

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
