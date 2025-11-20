package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lec7ral/WithWebSocket/internal/auth"
	"github.com/Lec7ral/WithWebSocket/internal/config"
	"github.com/Lec7ral/WithWebSocket/internal/logger"
	"github.com/Lec7ral/WithWebSocket/internal/repository"
	"github.com/Lec7ral/WithWebSocket/internal/websocket"
	"github.com/go-chi/chi/v5"
)

// Embed the 'static' directory.
//

//go:embed all:static
var embeddedFS embed.FS

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	appLogger := logger.New(cfg.LogLevel)
	slog.SetDefault(appLogger)

	ctx := context.Background()

	repo, err := repository.NewPostgresRepository(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to create repository", "error", err)
		os.Exit(1)
	}
	defer repo.Close()
	slog.Info("Database connection pool established.")

	authService := auth.NewService(cfg.JWTSecret)
	authHandler := auth.NewHandler(authService, repo)

	hub := websocket.NewHub(repo)
	go hub.Run()
	slog.Info("WebSocket Hub is running.")

	router := chi.NewRouter()
	wsHandler := websocket.NewHandler(hub, authService)

	// --- Static File Server Setup ---
	// Create a sub-filesystem that starts in the 'static' directory.
	staticFS, err := fs.Sub(embeddedFS, "static")
	if err != nil {
		slog.Error("Failed to create static filesystem", "error", err)
		os.Exit(1)
	}
	fileServer := http.FileServer(http.FS(staticFS))

	// Serve the frontend from the root path.
	router.Handle("/*", fileServer)

	// --- API and WebSocket Routes ---
	router.Post("/login", authHandler.HandleLogin)
	router.Route("/api", func(r chi.Router) {
		r.Get("/rooms", wsHandler.HandleGetRooms)
		r.Get("/users/{userID}", authHandler.HandleGetUser)
	})
	router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
			slog.Error("Failed to write health check response", "error", err)
		}
	})
	router.Get("/ws/{roomID}", wsHandler.ServeWS)

	// --- Graceful Shutdown Setup ---
	serverPort := ":" + cfg.ServerPort
	server := &http.Server{
		Addr:    serverPort,
		Handler: router,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Server starting", "port", serverPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	<-stopChan
	slog.Info("Shutdown signal received, starting graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown error", "error", err)
	}

	slog.Info("Server stopped gracefully.")
}
