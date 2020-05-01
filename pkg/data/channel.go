package data

import "time"

// Channel struct containing information about a joined channel
type Channel struct {
	ID            int
	ChannelID     string
	JoinedDate    time.Time
	PointsName    string
	CommandPrefix string
}
