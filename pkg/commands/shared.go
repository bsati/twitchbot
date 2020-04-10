package commands

import (
	"sync"
	"time"
)

// Command holds necessary information for invoking and metadata
type Command struct {
	UserCooldown   int
	GlobalCooldown int
	Description    string
	globalCooldown bool
	userCooldowns  map[string]bool
	cooldownMutex  *sync.RWMutex
}

// NewCommand creates a new command with given cooldown durations and descriptions
func NewCommand(userCooldown int, globalCooldown int, description string) *Command {
	return &Command{
		UserCooldown:   userCooldown,
		GlobalCooldown: globalCooldown,
		Description:    description,
		globalCooldown: false,
		userCooldowns:  make(map[string]bool),
		cooldownMutex:  &sync.RWMutex{},
	}
}

// HasCooldown checks whether the command has a cooldown for the specified user
func (c *Command) HasCooldown(user string) bool {
	c.cooldownMutex.RLock()
	defer c.cooldownMutex.RUnlock()
	if entry, ok := c.userCooldowns[user]; ok {
		return entry && c.globalCooldown
	}
	return c.globalCooldown
}

// AddCooldown adds a cooldown with specified duration for the user and globally
func (c *Command) AddCooldown(user string) {
	c.cooldownMutex.Lock()
	defer c.cooldownMutex.Unlock()
	if c.GlobalCooldown > 0 {
		c.globalCooldown = true
		time.AfterFunc(time.Duration(c.GlobalCooldown)*time.Second, func() {
			c.RemoveGlobalCooldown()
		})
	}
	if c.UserCooldown > 0 {
		c.userCooldowns[user] = true
		time.AfterFunc(time.Duration(c.UserCooldown)*time.Second, func() {
			c.RemoveCooldown(user)
		})
	}
}

// RemoveCooldown removes the specified user from the cooldown map
func (c *Command) RemoveCooldown(user string) {
	c.cooldownMutex.Lock()
	defer c.cooldownMutex.Unlock()
	delete(c.userCooldowns, user)
}

// RemoveGlobalCooldown removes the global cooldown
func (c *Command) RemoveGlobalCooldown() {
	c.cooldownMutex.Lock()
	defer c.cooldownMutex.Unlock()
	c.globalCooldown = false
}
