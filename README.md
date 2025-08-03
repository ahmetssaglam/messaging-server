# Messaging Server

## Overview
Messaging Server is a Go-based backend service designed to periodically process and deliver unsent messages. Every 2 minutes, a cron job fetches unsent messages from the database, sends them to a specified webhook, and then stores the processed records in a cache for fast retrieval.

## Features
- **Scheduled Cron Job:** Runs every 2 minutes to process unsent messages.
- **Database Integration:** Fetches unsent messages from a PostgreSQL database.
- **Webhook Delivery:** Sends messages to a configurable webhook endpoint.
- **Caching:** Stores sent message records in Redis for quick access.
- **Modular Structure:** Organized codebase for easy maintenance and scalability.

## Folder Overview
```
docker-compose.yaml      # Docker Compose configuration
Dockerfile               # Custom image build instructions
go.mod, go.sum           # Go dependencies
cmd/                     # Main application entry point
internal/                # Application logic and modules
pkg/                     # Utility packages
docs/                    # Swagger documentation
init/                    # SQL initialization scripts
```

## Project Structure
- **cmd/main.go:** Application entry point.
- **internal/cron:** Cron job logic.
- **internal/database:** Database access and models.
- **internal/handler:** Message and cron handlers.
- **internal/jobs:** Job logic for message processing.
- **internal/logging:** Logging utilities.
- **internal/models:** Data models.
- **internal/router:** API routing.
- **internal/configs:** Configuration management.
- **pkg/utils:** Utility functions.
- **docs/:** Swagger documentation.
- **init/:** SQL initialization scripts.

## Getting Started
### Prerequisites
- Go 1.24+ (for local development)
- Docker & Docker Compose
- PostgreSQL
- Redis

### Installation & Usage (Docker Compose)
1. Clone the repository:
```sh
git clone https://github.com/yourusername/messaging-server.git
cd messaging-server
```

2. Build and start the services using Docker Compose:
```sh
docker compose up -d --build
```
   This will build your custom image and start all required services (app, PostgreSQL, Redis).

3. The cron job will run automatically every 2 minutes, sending unsent messages from the database to the webhook and caching the results in Redis.

## API Documentation
- Swagger docs are generated in the `docs/` directory.
- To view Swagger UI, run the application and navigate to:
  ```
  http://localhost:8080/swagger/index.html
  ```

## Configurable Environment Variables

The application uses environment variables for configuration. You can set these in your Docker Compose file or a `.env` file. Key variables include:

| Variable Name         | Description                                  | Example Value                                                    |
|---------------------- |----------------------------------------------|------------------------------------------------------------------|
| APP_NAME              | Application name                             | Messaging Server V1                                              |
| LOG_LEVEL             | Logging level                                | DEBUG                                                           |
| MESSAGE_FETCH_LIMIT   | Number of messages to fetch per cron run     | 2                                                               |
| CRON_INTERVAL         | Cron job interval (in seconds)               | 120                                                             |
| MAX_CONCURRENT_JOBS   | Maximum number of concurrent jobs            | 5                                                               |
| SERVER_GRACE_PERIOD   | Grace period for server shutdown (seconds)   | 30                                                              |
| WEBHOOK_URL           | Webhook endpoint for message delivery        | https://webhook.site/ddd0be6e-d6c5-4859-a90f-9f72801bb182        |
| REDIS_HOST            | Redis host                                   | redis                                                           |
| REDIS_PORT            | Redis port                                   | 6379                                                            |
| REDIS_DB              | Redis database index                         | 0                                                               |
| REDIS_TTL             | Redis cache TTL (seconds)                    | 3600                                                            |
| POSTGRES_URI          | PostgreSQL connection URI                    | postgres://sample_user:sample_password@postgres_db:5432/messaging_db?sslmode=disable |

You can add or override these variables in your Docker Compose service definition under `environment:`.

## PostgreSQL Table Design

The main table used by the application is `messages`. Below is its schema:

| Column Name   | Type         | Constraints                | Description                  |
|-------------- |-------------|----------------------------|------------------------------|
| id            | VARCHAR(36) | PRIMARY KEY                | Unique message identifier    |
| content       | VARCHAR(255)| NOT NULL                   | Message content              |
| phone_number  | VARCHAR(20) | NOT NULL                   | Recipient phone number       |
| is_sent       | BOOLEAN     | NOT NULL, DEFAULT FALSE    | Message sent status          |

Sample rows are inserted for testing and development purposes.

## Cron Job Running Logic

The cron job is implemented in `internal/cron/cron.go` and works as follows:

- The `Cron` struct manages scheduled execution of a job function at a fixed interval (configured via environment variables).
- It uses a semaphore to limit the number of concurrent jobs, ensuring no more than `MAX_CONCURRENT_JOBS` run at the same time.
- When `Start()` is called, a goroutine is launched that triggers the job at each interval using a time ticker.
- For each tick, if concurrency limits allow, the job is executed in a new goroutine.
- The job function typically fetches unsent messages, sends them to the webhook, and updates their status.
- The cron can be stopped gracefully using a quit channel and WaitGroup.

This design ensures reliable, concurrent, and controlled execution of periodic tasks such as message delivery.

### Cleanup
To stop and remove all containers, networks, and volumes created by Docker Compose:
```sh
docker compose down --volumes
```
This will clean up your environment after testing or development.
