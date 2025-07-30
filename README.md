# Anti-Brute-Force Service

A service for protecting against brute-force attacks during user authentication. Provides an API for checking authentication attempts with rate limiting algorithm and managing white/black IP address lists.

## 🚀 Key Features

- **Rate Limiting**: Limiting authentication attempts by login, password, and IP address
- **White/Black Lists**: Managing lists of allowed and blocked IP subnets
- **gRPC API**: High-performance API for integration with other services
- **CLI Interface**: Command line interface for service administration
- **PostgreSQL**: Reliable storage for IP address list data

## 📋 Requirements

- Go 1.22+
- Docker & Docker Compose
- Task (task runner) - [installation](https://taskfile.dev/getting-started/)

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   gRPC Client   │    │  Anti-Brute-    │    │   PostgreSQL    │
│   (CLI/App)     │◄──►│  Force Service  │◄──►│   Database      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │  Rate Limiting  │
                       │  (Leaky Bucket) │
                       └─────────────────┘
```

## 🚀 Quick Start

### 1. Clone and Setup

```bash
git clone https://github.com/vagudza/anti-brute-force
cd anti-brute-force
```

### 2. Start Infrastructure

```bash
# Start PostgreSQL
task infra-up

# Apply database migrations
task migrations-up
```

### 3. Start Service

```bash
# Run in development mode
task run

# OR build and run in Docker
task docker-up
```

### 4. Install CLI

```bash
# Install CLI client
task install-cli
```

## ⚙️ Configuration

### Main Parameters

Configuration file: `config/app/config.local.yaml`

```yaml
limiters:
  login:
    maxAttemptsPerMinute: 10    # N - login attempts limit
    cleanupInterval: 1m
    ttl: 10m
  password:
    maxAttemptsPerMinute: 100   # M - password attempts limit
    cleanupInterval: 1m
    ttl: 10m
  ip:
    maxAttemptsPerMinute: 1000  # K - IP attempts limit
    cleanupInterval: 1m
    ttl: 10m

grpc:
  port: 13013

postgres:
  database: anti-brute-force
  username: anti-brute-force
  password: anti-brute-force
  host: localhost
  port: 5432
  sslmode: disable
```

### Environment Variables

- `CONFIG_PATH` - path to configuration file
- `ENV` - environment (development, production, ci)

## 🔧 Available Commands

### Infrastructure

```bash
# Start PostgreSQL
task infra-up

# Stop infrastructure
task infra-down

# Apply migrations
task migrations-up

# Rollback last migration
task migrations-down

# Create new migration
task migrations-new -- <migration-name>
```

### Development

```bash
# Format code
task fmt

# Lint code
task lint

# Format + lint
task pretty

# Generate proto files
task proto-gen
```

### Build and Run

```bash
# Build application
task build

# Run in development mode
task run

# Build and run in Docker
task docker-up

# Stop Docker containers
task docker-down
```

### Testing

```bash
# Unit tests
task test

# Integration tests
task integration-test

# CI checks (lint + tests + build)
task ci
```

### CLI

```bash
# Build CLI
task build-cli

# Install CLI
task install-cli
```

## 🧪 Testing

### Local Testing

```bash
# 1. Start infrastructure
task infra-up

# 2. Apply migrations
task migrations-up

# 3. Start service
task run

# 4. In another terminal - run tests
task test
task integration-test
```

### Integration Tests

Integration tests are located in the `/test` directory and verify:

- Authentication checking with rate limiting
- White/black list management
- Bucket reset functionality
- IP address and CIDR subnet validation

### CLI Usage Examples

```bash
# Add IP to whitelist
abf-cli whitelist add --subnet 192.168.1.0/24

# Add IP to blacklist
abf-cli blacklist add --subnet 10.0.0.0/8

# Reset bucket for login and IP
abf-cli bucket reset --login user@example.com --ip 192.168.1.100

# View lists
abf-cli whitelist list
abf-cli blacklist list
```

## 📡 API

### gRPC Endpoints

- `CheckAuth` - check authentication attempt
- `ResetBucket` - reset rate limiting bucket
- `AddToBlacklist/RemoveFromBlacklist` - manage blacklist
- `AddToWhitelist/RemoveFromWhitelist` - manage whitelist
- `GetBlacklist/GetWhitelist` - get lists
- `ClearBlacklist/ClearWhitelist` - clear lists

### Usage Example

```go
// Connect to service
conn, err := grpc.Dial("localhost:13013", grpc.WithInsecure())
client := pb.NewAntiBruteforceClient(conn)

// Check authentication
resp, err := client.CheckAuth(ctx, &pb.CheckAuthRequest{
    Login:    "user@example.com",
    Password: "password123",
    IP:       "192.168.1.100",
})
```

## 🐳 Docker

### Running in Docker

```bash
# Full stack (application + PostgreSQL)
task docker-up

# Infrastructure only
task infra-up

# Stop
task docker-down
```

### Docker Compose Profiles

- `infra` - PostgreSQL only
- `all` - full stack (application + PostgreSQL)

## 🔍 Monitoring

### Health Checks

- **Application**: `nc -z localhost 13013`
- **PostgreSQL**: `pg_isready -U anti-brute-force`

### Logs

Application logs are output to stdout/stderr and can be collected via Docker:

```bash
docker logs anti-brute-force
```

## 🛠️ Development

### Project Structure

```
├── api/proto/           # gRPC proto files
├── cmd/                 # Entry points
│   ├── app/            # Main application
│   └── abf-cli/        # CLI client
├── config/              # Configuration
├── internal/            # Internal packages
│   ├── app/            # Business logic
│   ├── bucket/         # Rate limiting
│   ├── config/         # Configuration
│   ├── iplist/         # White/Black lists
│   ├── storage/        # Database operations
│   └── transport/      # gRPC server
├── migrations/          # SQL migrations
├── pkg/cli/            # CLI package
└── test/               # Integration tests
```
