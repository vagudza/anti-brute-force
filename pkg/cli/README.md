# Anti-Bruteforce CLI

Command line interface for Anti-Bruteforce service.

## Installation

```bash
# From source code
go install github.com/vagudza/anti-brute-force/cmd/abf-cli@latest

# Or using task (from project root)
task install-cli

# Or just build without installation
task build-cli
```

## Usage

The CLI provides commands to manage buckets, whitelists, and blacklists for the Anti-Bruteforce service.

### Basic Usage

```bash
# Show help for all commands
abf-cli --help

# Show help for specific command
abf-cli bucket --help
abf-cli whitelist --help
abf-cli blacklist --help
```

### Bucket Management

The bucket command allows you to reset rate limiting counters for specific login/IP combinations.

```bash
# Reset bucket for specific login and IP
abf-cli bucket reset --login user@example.com --ip 192.168.1.100

# Reset all buckets for specific login
abf-cli bucket reset --login user@example.com

# Reset all buckets for specific IP
abf-cli bucket reset --ip 192.168.1.100
```

### Whitelist Management

The whitelist commands allow you to manage IP subnets that are always allowed to authenticate.

```bash
# Add subnet to whitelist
abf-cli whitelist add --subnet 192.168.1.0/24

# Remove subnet from whitelist
abf-cli whitelist remove --subnet 192.168.1.0/24

# List all whitelisted subnets
abf-cli whitelist list
```

### Blacklist Management

The blacklist commands allow you to manage IP subnets that are always blocked from authentication.

```bash
# Add subnet to blacklist
abf-cli blacklist add --subnet 10.0.0.0/8

# Remove subnet from blacklist
abf-cli blacklist remove --subnet 10.0.0.0/8

# List all blacklisted subnets
abf-cli blacklist list
```

## Configuration

The CLI can be configured using command line flags:

### Global Flags
- `--host` - Server host (default: localhost)
- `--port` - Server port (default: 13013)

You can use global flags with any command:

```bash
# Use custom host and port
abf-cli --host 192.168.1.10 --port 8080 whitelist list

# Multiple commands with same configuration
abf-cli --host 192.168.1.10 --port 8080 blacklist add --subnet 10.0.0.0/8
abf-cli --host 192.168.1.10 --port 8080 whitelist list
```

## Command Structure

The CLI uses a hierarchical command structure:

- `bucket` - Manage rate limiting buckets
  - `reset` - Reset buckets for specific login/IP
- `whitelist` - Manage whitelisted IP subnets
  - `add` - Add subnet to whitelist
  - `remove` - Remove subnet from whitelist
  - `list` - List all whitelisted subnets
- `blacklist` - Manage blacklisted IP subnets
  - `add` - Add subnet to blacklist
  - `remove` - Remove subnet from blacklist
  - `list` - List all blacklisted subnets

Each command has its own flags and options that can be viewed using the `--help` flag. 