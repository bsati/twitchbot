package data

// ChannelTextCommand is a basic command created by the channel owner / permitted users
type ChannelTextCommand struct {
	ID       int
	Aliases  []string
	Response string
}

// ChannelTimedCommand is a command thats' text is periodically sent to chat
type ChannelTimedCommand struct {
	ID    int
	Timer int
	Text  string
}

// ChannelModifiedCommand is a wrapper about an internal bot command that modifies the response text
type ChannelModifiedCommand struct {
	ID                int
	InternalCommandID int
	ModifiedText      string
}
