# twitchbot

Basic chat bot implementation for Twitch's IRC chat with command management, message and event handling.

## Usage

Run `go build cmd/bot.go` and then execute the resulting executable supplied with 2 arguments `NICK` and `OAUTH_TOKEN`.
Example usage on Windows: `.\bot.exe ExampleBot oauth:exampleToken`. For getting the OAuth-Token check [Twitch Reference](https://dev.twitch.tv/docs/irc/guide#connecting-to-twitch-irc)