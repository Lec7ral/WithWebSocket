package domain

// Message defines the structure for messages exchanged via WebSocket.
type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
	Sender  string `json:"sender,omitempty"`

	// RoomID is the identifier of the room this message belongs to.
	RoomID string `json:"room_id,omitempty"`
}
