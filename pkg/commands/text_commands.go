package commands

import (
	"github.com/bsati/twitchbot/pkg/twitchirc"
)

// BasicCommand interface for basic commands
type BasicCommand interface {
	Invoke([]string, *twitchirc.MessageEvent, *twitchirc.IRCClient)
}

// PingPong Command
type PingPong struct {
	cmd *Command
}

func newPingPong() *PingPong {
	return &PingPong{
		cmd: NewCommand(0, 0, "Ping pong!"),
	}
}

// Invoke BasicCommand interface implementation for PingPong
func (d *PingPong) Invoke(args []string, me *twitchirc.MessageEvent, irc *twitchirc.IRCClient) {
	irc.SendMessage(me.Channel, "Pong!")
}
