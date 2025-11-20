package domain

// DirectMessagePayload defines the structure for the payload of a direct message.
type DirectMessagePayload struct {
	// RecipientID is the ID of the user who should receive the message.
	RecipientID string `json:"recipient_id"`

	// Content is the actual text message.
	Content string `json:"content"`
}
