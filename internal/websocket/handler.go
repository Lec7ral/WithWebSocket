package websocket

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Lec7ral/WithWebSocket/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub         *Hub
	authService *auth.Service
}

func NewHandler(hub *Hub, authService *auth.Service) *Handler {
	return &Handler{
		hub:         hub,
		authService: authService,
	}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Token is required", http.StatusUnauthorized)
		return
	}

	claims, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		slog.Warn("Invalid WebSocket token received", "error", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade connection", "error", err)
		return
	}

	// The handler's job is just to validate and pass the request to the hub.
	regReq := &registrationRequest{
		claims: claims,
		conn:   conn,
		roomID: roomID,
	}
	h.hub.register <- regReq
}

// HandleGetRooms is the HTTP handler for the GET /api/rooms endpoint.
func (h *Handler) HandleGetRooms(w http.ResponseWriter, _ *http.Request) {
	rooms := h.hub.GetActiveRooms()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		slog.Error("Failed to write rooms response", "error", err)
	}
}
