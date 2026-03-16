# gamesync

gamesync is used to sync game save files to a server.

The project is seperated into client and server. 

The client is written in golang and syncs to the server using rsync.


## Features

- Sync game save files with server
- Multiple devices
- Multiple users

## Install client

### From the AUR
with yay:
```sh
yay -S gamesync
```

### Manual

Before attempting to install, make sure that you have the following dependencies installed:
- build dependencies
    - git
    - Go
    - Make
- package dependencies
    - OpenSSH
    - Rsync


#### Download the repo and cd into it:
```sh
git clone https://gitlab.com/PatrikBures/gamesync.git
cd gamesync
```

#### Checkout to the latest version tag (skip if you want to use the latest unstable version):
```sh
git checkout $(git tag -l 'v*.*.*' --sort=-version:refname | head -1)
```

##### If that by some reason did not work:
```sh
git tag -l
```
then run the following where VERSION is the latest version from the previous command
```sh
git checkout VERSION
```

#### compile and install:
```sh
make install PREFIX=$HOME/.local
```

> also make sure that `$HOME/.local/bin` is in your `$PATH`

#### Uninstall
while in the repo run:
```sh
make uninstall PREFIX=$HOME/.local
```


## Install server

### docker-compose with prebuilt image

make a docker-compose.yml file with this example content:
```yml
services:
  gamesync:
    image: patrikbures/gamesync:v2.6.0
    container_name: gamesync
    environment:
      - GAMESYNC_LOOP=20
    ports:
      - "2828:22"
    volumes:
      - ./config:/config
      - ./data:/data
      - gamesync_host_keys/host_keys:/etc/ssh/keys
      - /etc/localtime:/etc/localtime:ro
volumes:
  gamesync_host_keys:
```

> If there is a new version use that as this example config will probably not be updated for every new version

then run
```sh
docker compose up -d
```

### build your own image with and run with docker-compose

clone the repo and cd into it and run:
```sh
make build-state
docker compose up -d --build
```

## Set up user on server

- create a file at `config/users/USER` where USER is the new name of your user
- in that new file add the public ssh keys that user will use

example content of a user file:
```
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILhrZ64KRubJJlHMcJ6ckbQx5XU2yh6/+AdBPTSxvluh user@device1
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIORjJGTNJucZcWfsWdl71LqYE/w2sxv2I3xuVwpFvcHx user@device2
```
make sure that the clients use the same USER name

## Set up client config


> When gamesync is ran it will create the dir but not the config.yml which needs to be made manually

### Create and configure ssh

You can create an ssh key with this command
```sh
ssh-keygen -f ~/.ssh/gamesync -N ''
```

### Add config for gamesync in ssh

example content of `~/.ssh/config`:
```
Host gamesync
    Hostname 192.168.1.10
    User username
    Port 2828
    IdentityFile ~/.ssh/gamesync
    ControlMaster auto
    ControlPath ~/.ssh/sockets/%r@%h-%p
    ControlPersist 2m
```

> The __Control__ config options make a persistent connection which lasts for 2 minutes. 
> This makes commands which run multiple ssh commands one after the other significantly faster speeding up most commands.


### Configure client

#### Location of config

the config file is read from
```sh
$XDG_CONFIG_HOME/gamesync/config.yml
```
or
```sh
$HOME/.config/gamesync/config.yml
```
depending if `XDG_CONFIG_HOME` is set

On a mac, it should be located at
```sh
$HOME/Library/Application/gamesync/config.yml
```

#### With ssh config (recommended)

This is if you set up an ssh config in the previous step

example config.yml:
```yml
server:
  ssh_host: gamesync
```

#### Without ssh config

> This might in the future become removed so it is recommended to not use it
> It also does not have the __Control__ options mentioned earlier.

example config.yml:
```yml
server:
  host: 192.168.0.10
  user: user1
  port: 2828
  identity_file: /home/user1/.ssh/gamesync
```


## Add games to config

### Manual way

You can manually add the games to the config by modifying the config.yml file.

example content of config.yml games:

```yml
games:
- id: game_id_1
  save_path: /full path/to/game 1
- id: game_id_2
  save_path: /full path/to/game 2
```

the game id is the id you will use in other commands to sync your save files, like in the wrap command explained later.

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


Relative paths can also be used
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
Replace __NAME__ with the name of the app or game.

2. Then copy that file to your users application directory. Usually located at `~/.local/share/applications`

If your desktop file already is somewhere in your home dir, you should still create a copy of it with a suffix. 
Just in case you regret your modifications.

2. Open the copied .desktop file and modify the Exec command to something like this
```ini
Exec=gamesync wrap GAME_ID -- [whatever was originally here]
```

### Useful wrap flags

- `-n`: Sends notification whenever it pushed or pulled. Notifying you if it succeded or failed.
- `--no-pull`: Useful if you want to just try out the wrap command or use as a backup without modifying your local save files.
- `-e`: Exits on error. For example, if the pull failed it exits, preventing the game from launching.

## Client environmental variables

- `GAMESYNC_CONFIG`: sets the config file to use, does not overwrite --config flag
- `GAMESYNC_STATE`: sets the state dir where sync states are stored

## Container environmental variables

- `GAMESYNC_LOOP`
    - when set to a number above 0 it will create users that often in seconds
    - if unset or below 1, only creates user on startup
    - useful if you want to add or remove users without restarting container

## Copyright

Copyright (C) 2025-2026 Patrik Bures
