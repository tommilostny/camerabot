# camerabot

[![Build Status](https://travis-ci.org/cooldarkdryplace/camerabot.svg?branch=master)](https://travis-ci.org/cooldarkdryplace/camerabot)
[![Go Report Card](https://goreportcard.com/badge/github.com/cooldarkdryplace/camerabot)](https://goreportcard.com/report/github.com/cooldarkdryplace/camerabot)
[![GoDoc](https://godoc.org/github.com/cooldarkdryplace/camerabot?status.svg)](https://godoc.org/github.com/cooldarkdryplace/camerabot)

## Building

You need to get sources and cross-compile them for ARM architecture. These can be easily done by these commands, assuming you have Go installed.

1. `go get github.com/tommilostny/camerabot`
2. `env GOOS=linux GOARCH=arm go build -v github.com/tommilostny/camerabot/cmd/camerabot`

As a result you will have binary suitable for running on Raspberry Pi. Copy it on device and proceed.

## Software
Telegram bot that makes a photo and sends it to chat. 

I use this bot to monitor kiln temperature and make sure workshop is not on fire yet.
Go part is responsible for interacting with Telegram API. Application uses long polling because in my case device is located behind two NATs. 
Uses `raspistill` (via `os.exec`) to make photos.
Parametrized commands for raspistill are stored in external bash scripts.

### Running bot
1. Setup Raspberry Pi and Pi camera.
2. Set environment variable `TOKEN` with your bot token (Botfather can provide you with the one).
3. Set environment variable `ALLOWED_CHAT_IDS` with chat IDs you want to have access to the camera (separated by ;).
4. Use systemd config to start as a service or simply run the app from the console.
5. Start direct conversation with bot or add bot to group chat if you are interested in broadcasting your kiln paranoia.

### Commands
1. `/pic` sends ordinary photo.
2. `/zoom` sends zoomed and croped region of interest. Kiln controller in my case.

---

## Use case after fork
Apartment monitor camera that reports via Telegram Bot to a few roommates (``AllowedChatIDs``).

This use case was enhanced by adding [gocron](https://github.com/go-co-op/gocron) scheduler, that sends a photo from the camera to all allowed Chat IDs every **30 minutes** using the ``/pic`` command handler.

![Apartment monitor](img/5782978610443958432_121.jpg)

---

#### Examples
![pic processing result](https://cloud.githubusercontent.com/assets/6103939/23331112/898d1204-fb67-11e6-8285-6efc5ba7816b.png)
![another pic result](https://cloud.githubusercontent.com/assets/6103939/23331113/92065df0-fb67-11e6-9d0f-d8adc245f9a3.png)
![zoom result](https://cloud.githubusercontent.com/assets/6103939/23331114/9b4fa8e4-fb67-11e6-876e-318642f38dfc.png)

## Hardware
Currently runs on a Raspberry Pi. Using onboard V2 camera.
