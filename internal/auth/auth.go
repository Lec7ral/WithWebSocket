package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Lec7ral/WithWebSocket/internal/domain"
	"github.com/Lec7ral/WithWebSocket/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid" // Import the missing package
)

// Claims defines the structure of the JWT claims.
type Claims struct {
	UserID    string `json:"userId"`
	Username  string `json:"username"`
	SessionID string `json:"sessionId"`
	jwt.RegisteredClaims
}

// Service provides authentication-related operations.
type Service struct {
	jwtSecret  string
	jwtExpires time.Duration
}

// NewService creates a new authentication service.
func NewService(jwtSecret string) *Service {
	return &Service{
		jwtSecret:  jwtSecret,
		jwtExpires: 24 * time.Hour,
	}
}

// GenerateToken generates a new JWT for a given user.
func (s *Service) GenerateToken(user *domain.User) (string, error) {
	claims := &Claims{
		UserID:    user.ID,
		Username:  user.UserName,
		SessionID: uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken validates a JWT string and returns the claims if valid.
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Handler handles HTTP requests for authentication.
type Handler struct {
	service *Service
	repo    repository.Repository
}

// NewHandler creates a new authentication handler.
func NewHandler(service *Service, repo repository.Repository) *Handler {
	return &Handler{
		service: service,
		repo:    repo,
	}
}

// LoginRequest defines the structure of a login request body.
type LoginRequest struct {
	Username string `json:"username"`
}

// LoginResponse defines the structure of a successful login response.
type LoginResponse struct {
	Token string `json:"token"`
}

// HandleLogin handles the /login endpoint.
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	user, err := h.repo.FindOrCreateUserByUsername(context.Background(), req.Username)
	if err != nil {
		slog.Error("Failed to find or create user", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := h.service.GenerateToken(user)
	if err != nil {
		slog.Error("Failed to generate token", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(LoginResponse{Token: token}); err != nil {
		slog.Error("Failed to write login response", "error", err)
	}
}

// HandleGetUser handles requests to get a user's public information.
func (h *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	user, err := h.repo.FindUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	publicUser := domain.User{
		ID:       user.ID,
		UserName: user.UserName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(publicUser); err != nil {
		slog.Error("Failed to write user response", "error", err)
	}
}
