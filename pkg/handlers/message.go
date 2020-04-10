package handlers

import (
	"github.com/bsati/twitchbot/pkg/commands"
	"github.com/bsati/twitchbot/pkg/core"
	"github.com/bsati/twitchbot/pkg/twitchirc"
)

// MessageHandler struct that wraps message handling funcs
type MessageHandler struct {
	*core.Env
}

// GetCommandHandler returns the func for message handling that parses commands
func (mh *MessageHandler) GetCommandHandler() func(irc *twitchirc.IRCClient, me *twitchirc.MessageEvent) {
	return func(irc *twitchirc.IRCClient, me *twitchirc.MessageEvent) {
		cmd, args := mh.Env.CommandMatcher.Match(me.Message)
		if cmd != nil {
			switch v := cmd.(type) {
			case commands.BasicCommand:
				v.Invoke(args, me, irc)
			}
		}
	}
}
