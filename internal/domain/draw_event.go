package domain

// DrawEventPayload defines the structure for the data associated with a single drawing event.
type DrawEventPayload struct {
	// The X coordinate of the event.
	X int `json:"x"`

	// The Y coordinate of the event.
	Y int `json:"y"`

	// Optional: The color of the stroke.
	Color string `json:"color,omitempty"`

	// Optional: The width of the stroke.
	LineWidth int `json:"lineWidth,omitempty"`
}
