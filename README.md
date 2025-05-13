# voidkit-go

A Go boilerplate project with Docker setup for development.

## Prerequisites

- Go 1.23.3 or later
- Docker and Docker Compose (for development)
- Make (optional, for using Makefile commands)

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/Gambitier/voidkit-go.git
cd voidkit-go
```

2. Set up environment:

Copy & edit docker.env to customize your settings

```bash
# Copy the sample environment file
cp docker.sample.env docker.env
```

3. Start the development environment:

```bash
# Start all services (Redis and PostgreSQL)
make up

# Run the Go server
make serve

# Or do both with
make dev
```

## Available Commands

- `make up`: Start Docker services (Redis and PostgreSQL)
- `make down`: Stop Docker services
- `make serve`: Run the Go server locally
- `make dev`: Start services and run the server (up + serve)
- `make build`: Build the Go binary
- `make test`: Run specific test (use TEST=TestName)
- `make tests`: Run all tests
- `make clean`: Clean up build artifacts and stop services
- `make proto`: Generate protobuf and gRPC code
- `make logs`: View service logs

## Development

The project uses:
- Redis for caching
- PostgreSQL for data storage
- Docker for service containerization
- Make for common development tasks

## Project Structure

```
voidkit-go/
├── cmd/                   # Application entry points
├── internal/              # Private application code
├── pkg/                   # Public library code
├── tests/                 # Test files
├── docker-compose.dev.yml # Development Docker configuration
├── docker.env             # Docker environment variables
├── docker.sample.env      # Sample Docker environment file
└── Makefile               # Development commands
```
