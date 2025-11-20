package domain

// RoomState represents the complete state of a room at a given moment.
// It will be sent to a user when they join the room.
type RoomState struct {
	// Users is the list of all users currently in the room.
	Users []*User `json:"users"`

	Whiteboard *WhiteboardState `json:"whiteboard"`
}
