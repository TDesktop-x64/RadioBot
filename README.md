# Telegram Radio Controller (For TDesktop-x64 only)

[![GitHub release](https://img.shields.io/github/v/release/c0re100/RadioBot.svg)](https://github.com/c0re100/RadioBot/releases/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

An experimental bot for controlling your TDesktop-x64(with Radio mode) music streaming session.

[TDesktop-x64](https://github.com/TDesktop-x64) Radio mode is an experimental feature for streaming music through voice
chat without audio filtering/processing.

That's mean you can stream your foobar2000/DeaDBeeF music library in TDesktop-x64.

## Features

Basically, request, skip and search a song are supported now.

### Commands

#### User&Admin:

* `/request` Request a song with inline button menu.
* `/skip` Start a poll to skip current song.
* `/search <name>` or `/nom <name>` Search a song with track name or artist.

#### Admin

* `/play` Play a song
* `/stop` Stop a song
* `/pause` Pause a song
* `/random` Skip to next random song
* `/reload` Reload songs list and controller config

### [WIP] Userbot mode

Unluckily, user radio bot IS NOT 100% implemented, I hope I will finish in someday...also help wanted!!!

If you're familiar with [pion/webrtc](https://github.com/pion/webrtc), please feel free to contact me
via [TDesktop-x64 Chat](https://t.me/tg_x64_chat) :)

## Quick Start

### Requirements

1. [Telegram Desktop x64](https://github.com/TDesktop-x64/tdesktop) with Radio mode
2. [foobar2000](https://www.foobar2000.org/) or [DeaDBeeF](https://deadbeef.sourceforge.io/)
3. Virtual Audio (Windows: _VB-Cable_ Linux: _PulseAudio_ macOS: _BlackHole_)
4. [Beefweb](https://github.com/hyperblast/beefweb)

*If you can't access [Beefweb Web Interface with default port](http://localhost:8880), please change the port through [beefweb plugins settings](images/beefweb_port.png) and edit [`beefweb_port`](#configuration) in config.json.*

### Setup

1. Install virtual audio
2. Reboot if needed
3. Copy config.json.sample to config.json
4. Edit config.json (See [Configuration](#configuration))
5. Open foobar2000 or DeaDBeeF, set virtual audio as your music player output
6. Open TDesktop-x64, set virtual audio as your microphone
7. Add your control bot to group and join a voice chat to play your song.
8. Done~

### Building

[tdlib](https://github.com/tdlib/td#building)

```
git clone https://github.com/c0re100/RadioBot
cd RadioBot
go build
```

### Prebuilt

[Release](https://github.com/c0re100/RadioBot/releases)

### Configuration

| Parameter        | Type    | Description                                                       |
| ---------------- | ------- | ----------------------------------------------------------------- |
| `api_id`         | String  | Obtain API ID from [my.telegram.org](https://my.telegram.org)     |
| `api_hash`       | String  | Obtain API Hash from [my.telegram.org](https://my.telegram.org)   |
| `bot_token`      | String  | Obtain bot token from [@BotFather](https://t.me/BotFather)        |
| `chat_id`        | Integer | Chat identifier, It can be empty if your group is public.         |
| `chat_username`  | String  | Empty if your group is private and fill the `chat_id` field.      |
| `pinned_message` | Integer | Info message of current song playing                              |
| `beefweb_port`   | Integer | Beefweb Port, Default: 8880                                       |
| `playlist_id`    | String  | foobar2000/DeaDBeeF playlist identifier,                          |
|                  |         | Obtain ID from [Beefweb API](http://localhost:8880/api/playlists) |

| **Limit**                 | Type    | Description                                            |
| ------------------------- | ------- | ------------------------------------------------------ |
| `chat_select_limit`       | Integer | Select page of rate limit for Chat, Default: 5         |
| `private_select_limit`    | Integer | Select page of rate limit for Private, Default: 10     |
| `row_limit`               | Integer | Number of rows, Default: 10                            |
| `queue_limit`             | Integer | Max queue songs, Default: 50                           |
| `recent_limit`            | Integer | Max recent songs, Default: 50                          |
| `request_song_per_minute` | Integer | Request a songs per minute limit, Default: 1 minute(s) |

| **Vote**             | Type    | Description                                                    |
| -------------------- | ------- | -------------------------------------------------------------- |
| `enable`             | Boolean | If true, user can start a vote to skip current song.           |
|                      |         | Default: true                                                  |
| `vote_time`          | Integer | Vote time, Default: 45s                                        |
| `update_time`        | Integer | Update vote status each n seconds, Default: 15s                |
| `release_time`       | Integer | Lock the vote n seconds after vote ended, Default: 600s        |
| `percent_of_success` | Float64 | Success percentage, Default: 40%                               |
| `participants_only`  | Boolean | If true, only participants which are in a voice chat can vote! |
|                      |         | Default: true                                                  |
| `user_must_join`     | Boolean | If true, Only users which are in the group can vote.           |
|                      |         | Default: false                                                 |

| **Web**  | Type    | Description                                               |
| -------- | ------- | --------------------------------------------------------- |
| `enable` | Boolean | If false, switch to Userbot mode.                         |
|          |         | If `participants_only` is true, please enable web server, |
|          |         | and fill the Radio controller url in TDesktop-x64.        |
|          |         | Default: true, since userbot is not finished.             |
| `port`   | Integer | Server Port, Default: 2468                                |