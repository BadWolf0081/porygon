# porygon
Porygon is a fork of https://github.com/roundaboutluke/porygon : Discordopole for Golbat, written in go.

<img src="https://i.imgur.com/Q7jKuVY.png" width="150" title="hover text">

# Requirements

[go 1.21](https://go.dev/doc/install)

# Installation

1. Git clone the repo `git clone https://github.com/roundaboutluke/porygon.git`
2. `cp default.toml config.toml` & adjust config.toml accordingly
3. `go build .`
5. `pm2 start ./porygon --name porygon`

---

There's also a Docker image `ghcr.io/lenisko/porygon:main`


# Further Customisation

Basic customisation of Porygon's localisation and layout. Simply `cp templates/current.json templates/current.override.json` and edit accordingly, using the examples in current.json. Note that some of the more generic emojis are contained within this.

To get emoji working, you want to create two guilds where one would be for bot to join and store emojis. On emote guild execute `/create-emotes` and `/list-emotes` to upload and list emotes for template needs. Emotes can be used on all guilds bot is logged.

# Updating

1. `git pull`
3. `go build .`
3. `pm2 restart porygon`

# Discord Permissions

Porygon requires your bot have the server permissions **Send Messages**, **Read Message History** and **Embed Links** to function.
