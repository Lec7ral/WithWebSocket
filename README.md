# CollabSphere

**CollabSphere** is a real-time collaboration application backend built with Go. It provides a robust foundation for features like multi-room chat, direct messaging, and a persistent collaborative whiteboard, all powered by WebSockets and a modern Go technology stack.

This project serves as a comprehensive example of building a stateful, concurrent, and production-ready real-time server.

## Features

- **Real-Time Communication**: Low-latency bidirectional communication using WebSockets.
- **Multi-Room Chat**: Clients can join different rooms, and communication is isolated between them.
- **Direct Messaging**: Users can send private messages to each other.
- **Collaborative Whiteboard**: A persistent, real-time whiteboard where drawing events are broadcast to all room participants.
- **Stateful Backend**: The server maintains the state of rooms, including user lists and whiteboard content.
- **Authentication**: Secure user authentication using JSON Web Tokens (JWT).
- **Persistence**: Chat history and whiteboard states are saved to a PostgreSQL database.
- **Resilience**:
    - **Graceful Shutdown**: The server shuts down gracefully, ensuring no data is lost.
    - **Rate Limiting**: Protects the server from being overwhelmed by limiting the number of messages a client can send.
- **Professional Tooling**:
    - **Structured Logging**: All logs are in JSON format using `slog` for better observability.
    - **Centralized Configuration**: Flexible configuration management using Viper (reads from `config.yaml` and environment variables).
    - **Embedded Frontend**: A functional test client is embedded directly into the Go binary for easy testing and deployment.

## Tech Stack & Architecture

The application follows a clean, modular architecture. The core logic is centered around a concurrent `Hub` that manages clients and message routing.

- **Language**: Go
- **Web Framework**: [Chi](https://github.com/go-chi/chi) (for HTTP routing)
- **WebSockets**: [Gorilla WebSocket](https://github.com/gorilla/websocket)
- **Database**: PostgreSQL
- **Database Driver**: [pgx](https://github.com/jackc/pgx)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Logging**: `slog` (Go 1.21+ structured logging)
- **Authentication**: [jwt-go](https://github.com/golang-jwt/jwt)
- **UUIDs**: [google/uuid](https://github.com/google/uuid)

## Getting Started

Follow these instructions to get the CollabSphere server running on your local machine.

### 1. Prerequisites

- **Go**: Version 1.21 or later.
- **PostgreSQL**: A running instance of PostgreSQL.

### 2. Database Setup

You need to create a database and the necessary tables for the application.

1.  **Create a Database**: Connect to your PostgreSQL instance and create a new database. You can name it `collabsphere_db` or similar.

    ```sql
    CREATE DATABASE collabsphere_db;
    ```

2.  **Create Tables**: Execute the `schema.sql` file located in the project root against your newly created database. This will create the `users`, `messages`, and `whiteboards` tables.

    ```bash
    # Example using psql
    psql -d collabsphere_db -f schema.sql
    ```

### 3. Application Configuration

The application uses the `config.yaml` file for default configuration, which can be overridden by environment variables.

1.  **Copy the Configuration**: If you haven't already, ensure you have a `config.yaml` file in the root of the project.

2.  **Update Database URL**: The most important step is to update the database connection string. You can either:
    -   **Edit `config.yaml`**: Change the `database.url` to match your PostgreSQL credentials.
        ```yaml
        database:
          url: "postgres://YOUR_USER:YOUR_PASSWORD@localhost:5432/collabsphere_db?sslmode=disable"
        ```
    -   **(Recommended) Use an Environment Variable**: Set the `DATABASE_URL` environment variable. Viper will automatically pick it up.
        ```bash
        # On Linux/macOS
        export DATABASE_URL="postgres://YOUR_USER:YOUR_PASSWORD@localhost:5432/collabsphere_db?sslmode=disable"

        # On Windows (PowerShell)
        $env:DATABASE_URL="postgres://YOUR_USER:YOUR_PASSWORD@localhost:5432/collabsphere_db?sslmode=disable"
        ```

### 4. Run the Application

With the database and configuration ready, you can now run the server.

1.  **Install Dependencies**: Open a terminal in the project root and download the Go modules.
    ```bash
    go mod tidy
    ```

2.  **Run the Server**:
    ```bash
    go run ./cmd/collabsphere
    ```

3.  **Access the Application**: The server will start, and you will see structured logs in your terminal.
    -   Open your web browser and navigate to `http://localhost:8080`.
    -   You will be greeted with the embedded test client, ready to use.

You have now successfully set up and launched the CollabSphere application!
