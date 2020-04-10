package commands

// CommandListEntry struct holding info for a command in the commandlist
type CommandListEntry struct {
	Aliases []string
	Command interface{}
}

// BuildCommandList returns a map containing all currently implemented commands
func BuildCommandList() map[int]CommandListEntry {
	result := make(map[int]CommandListEntry)
	result[0] = *makeCommandListEntry(newPingPong(), "ping")
	return result
}

func makeCommandListEntry(cmd interface{}, aliases ...string) *CommandListEntry {
	return &CommandListEntry{
		Aliases: aliases,
		Command: cmd,
	}
}
