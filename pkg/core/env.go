package core

import (
	"sync"

	"github.com/bsati/twitchbot/pkg/commands"
)

// Env holds environment variables
type Env struct {
	Config
	*commands.CommandMatcher
	ChannelCommandMatcherCache      map[string]*commands.CommandMatcher
	ChannelCommandMatcherCacheMutex *sync.RWMutex
}
