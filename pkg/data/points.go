package data

// UserPoints struct containing information about the points of a user in a certain channel
type UserPoints struct {
	ID       int
	Username string
	Amount   int
	// Foreign relation
	ChannelID string
}

// GetPoints returns the UserPoints struct read from the db if one exists for given username and channel
func GetPoints(Username string, ChannelID string) (UserPoints, error) {

}

// CreatePoints creates a new db entry for the given username and channel
func CreatePoints(Username string, ChannelID string) error {

}

// UpdatePoints updates an existing points db entry with the new amount
func UpdatePoints(points UserPoints, NewAmount int) error {

}
