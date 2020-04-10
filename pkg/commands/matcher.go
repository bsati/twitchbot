package commands

import (
	"strings"
	"sync"
)

// CommandMatcher holds information dealing with matching inputs to commands
type CommandMatcher struct {
	prefix   string
	sep      string
	mutex    *sync.RWMutex
	commands map[string]interface{}
}

// NewCommandMatcher returns a pointer to a new CommandMatcher with specified seperator and prefix
func NewCommandMatcher(prefix string, seperator string) *CommandMatcher {
	return &CommandMatcher{
		prefix:   prefix,
		sep:      seperator,
		mutex:    &sync.RWMutex{},
		commands: make(map[string]interface{}),
	}
}

// Register adds a command with given aliases to the command map
func (cm *CommandMatcher) Register(aliases []string, command interface{}) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	for _, alias := range aliases {
		cm.commands[alias] = command
	}
}

// Deregister deletes entries with given aliases
func (cm *CommandMatcher) Deregister(aliases ...string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	for _, alias := range aliases {
		delete(cm.commands, alias)
	}
}

// Match looks for a command with the alias of the first input word if the prefix matches
func (cm *CommandMatcher) Match(text string) (interface{}, []string) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	split := strings.Split(text, cm.sep)
	if string(split[0][0]) == cm.prefix {
		if cmd, ok := cm.commands[split[0][1:]]; ok {
			return cmd, split[1:]
		}
	}
	return nil, nil
}
