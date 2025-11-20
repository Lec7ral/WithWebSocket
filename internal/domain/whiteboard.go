package domain

type DrawEvent struct {
	Type    string           `json:"type"`
	Payload DrawEventPayload `json:"payload"`
}
type WhiteboardState struct {
	Events []DrawEvent `json:"events"`
}
