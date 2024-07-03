# Effective Mobile Test Task

This is a test task for Effective Mobile.

## Running Locally

1. **Ensure PostgreSQL is installed**: Make sure PostgreSQL is installed and running on your machine.

2. **Run with `go run`**:

   - Navigate to your project directory in the terminal.
   - Use `go run ./cmd` to run your Go application. This assumes your main package is located under `cmd`.

   The application will be accessible locally.

## Running with Docker Compose

1. **Build and Start Docker Containers**:

   - Ensure Docker is installed on your machine.
   - Use the following commands:
     ```
     docker-compose build  # Build the Docker containers
     docker-compose up     # Start the application
     ```

   The application will be accessible locally.

## Swagger Documentation

- Swagger documentation can be accessed at `docs/swagger/index.html`.

## Default Port

By default, the application is accessible on port 8080.

## SQL Schema and Mock Files

- SQL schema and mock files can be found at `pkg/migrations/sql/`.
