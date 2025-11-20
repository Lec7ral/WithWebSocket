package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Lec7ral/WithWebSocket/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for database operations.
type Repository interface {
	SaveMessage(ctx context.Context, msg *domain.Message) error
	GetMessagesByRoom(ctx context.Context, roomID string, limit int) ([]*domain.Message, error)
	FindOrCreateUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetWhiteboardState(ctx context.Context, roomID string) (*domain.WhiteboardState, error)
	SaveWhiteboardState(ctx context.Context, roomID string, state *domain.WhiteboardState) error
	FindUserByID(ctx context.Context, userID string) (*domain.User, error) // New method
	Close()
}

// PostgresRepository is the PostgreSQL implementation of the Repository.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new repository and connects to the database.
func NewPostgresRepository(ctx context.Context, connString string) (*PostgresRepository, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	return &PostgresRepository{pool: pool}, nil
}

// Close closes the database connection pool.
func (r *PostgresRepository) Close() {
	r.pool.Close()
}

// SaveMessage saves a message to the database.
func (r *PostgresRepository) SaveMessage(ctx context.Context, msg *domain.Message) error {
	if msg.Type != "text_message" {
		return nil
	}
	payloadStr, ok := msg.Payload.(string)
	if !ok {
		return fmt.Errorf("invalid payload type for text_message")
	}
	query := `INSERT INTO messages (room_id, sender_id, payload) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, msg.RoomID, msg.Sender, payloadStr)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	return nil
}

// GetMessagesByRoom retrieves the last N messages for a given room.
func (r *PostgresRepository) GetMessagesByRoom(ctx context.Context, roomID string, limit int) ([]*domain.Message, error) {
	query := `
		SELECT room_id, sender_id, payload, timestamp
		FROM messages
		WHERE room_id = $1
		ORDER BY timestamp DESC
		LIMIT $2`

	rows, err := r.pool.Query(ctx, query, roomID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var msg domain.Message
		var timestamp time.Time
		if err := rows.Scan(&msg.RoomID, &msg.Sender, &msg.Payload, &timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan message row: %w", err)
		}
		msg.Type = "archived_text_message"
		messages = append(messages, &msg)
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// FindOrCreateUserByUsername finds a user by username or creates a new one.
func (r *PostgresRepository) FindOrCreateUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	queryFind := `SELECT id, username FROM users WHERE username = $1`
	var user domain.User
	err := r.pool.QueryRow(ctx, queryFind, username).Scan(&user.ID, &user.UserName)

	if err == nil {
		return &user, nil
	}

	if err == pgx.ErrNoRows {
		newID := uuid.NewString()
		queryCreate := `INSERT INTO users (id, username) VALUES ($1, $2) RETURNING id, username`
		var newUser domain.User
		err := r.pool.QueryRow(ctx, queryCreate, newID, username).Scan(&newUser.ID, &newUser.UserName)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
		return &newUser, nil
	}

	return nil, fmt.Errorf("failed to find user: %w", err)
}

// GetWhiteboardState retrieves the current state of the whiteboard for a given room.
func (r *PostgresRepository) GetWhiteboardState(ctx context.Context, roomID string) (*domain.WhiteboardState, error) {
	query := `SELECT state FROM whiteboards WHERE room_id = $1`
	var stateJSON []byte
	err := r.pool.QueryRow(ctx, query, roomID).Scan(&stateJSON)

	if err == pgx.ErrNoRows {
		return &domain.WhiteboardState{Events: []domain.DrawEvent{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query whiteboard state: %w", err)
	}

	var state domain.WhiteboardState
	if err := json.Unmarshal(stateJSON, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal whiteboard state: %w", err)
	}

	return &state, nil
}

// SaveWhiteboardState saves or updates the state of the whiteboard for a given room.
func (r *PostgresRepository) SaveWhiteboardState(ctx context.Context, roomID string, state *domain.WhiteboardState) error {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal whiteboard state: %w", err)
	}

	query := `
		INSERT INTO whiteboards (room_id, state, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (room_id) DO UPDATE
		SET state = EXCLUDED.state, updated_at = NOW()`

	_, err = r.pool.Exec(ctx, query, roomID, stateJSON)
	if err != nil {
		return fmt.Errorf("failed to save whiteboard state: %w", err)
	}

	return nil
}

// FindUserByID finds a single user by their ID.
func (r *PostgresRepository) FindUserByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `SELECT id, username FROM users WHERE id = $1`
	var user domain.User
	err := r.pool.QueryRow(ctx, query, userID).Scan(&user.ID, &user.UserName)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user with ID %s not found", userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return &user, nil
}
