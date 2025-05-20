# Delivr

A simple Go application that executes commands (Docker, Git, or any other CLI tool) and sends notifications to Discord.

> This project was built with [Windsurf](https://windsurf.com) and [Claude](https://www.anthropic.com/claude) AI assistance.

## Features

- Execute various types of commands (Docker, Git, etc.)
- Send command output to Discord via webhook
- Log command outputs to files with automatic rotation
- Configure via JSON or YAML with support for custom paths
- Flexible configuration with optional sections
- Error handling and comprehensive reporting
- One-time execution or daemon mode
- Automated builds with GitHub Actions

## Installation

```bash
# Clone the repository
git clone https://github.com/ndious/delivr.git
cd delivr

# Install dependencies
go mod tidy

# Build the application
go build
```

## Quick Start

Generate a default configuration file:

```bash
./delivr --init
```

Edit the generated configuration file `.delivr.yml` avec your Discord webhook URL and customize your commands.

## Usage

```bash
# Run with default config file (.delivr.yml or .delivr.json in current directory)
./delivr

# Specify a custom config file
./delivr --config /path/to/your/config.json

# Run in daemon mode (doesn't exit after executing all commands)
./delivr --daemon

# Generate a default configuration file
./delivr --init

# Generate a configuration file at a specific location
./delivr --init --out /path/to/new/.delivr.yml
```

## Configuration

The configuration file supports both JSON and YAML formats. The application will look for configuration files in the following order:

1. `.delivr.yml` in the current directory
2. `.delivr.json` in the current directory
3. `config.yml` in the current directory (deprecated)
4. `config.json` in the current directory (deprecated)
5. `config.yml` in the user's home directory under `.delivr/`
6. `config.json` in the user's home directory under `.delivr/`

### YAML Configuration Example (Recommended)

```yaml
# Optional working directory for commands
workingDir: /path/to/your/working/directory

# Optional Docker configuration
docker:
  host: unix:///var/run/docker.sock

# Optional logging configuration
logs:
  directory: ./logs
  maxSize: 10
  maxAge: 30
  maxBackups: 5
  compress: true

# Discord webhook configuration (required)
discord:
  channelId: https://discord.com/api/webhooks/YOUR_WEBHOOK_URL

# Commands to execute (required)
commands:
  - name: Show Docker Status
    description: Lists all running Docker containers
    command: docker
    args: [ps, -a]
  
  - name: Git Status
    description: Shows the working tree status
    command: git
    args: [status]
```

### JSON Configuration Example

```json
{
  "workingDir": "/path/to/your/working/directory",
  "docker": {
    "host": "unix:///var/run/docker.sock"
  },
  "logs": {
    "directory": "./logs",
    "maxSize": 10,
    "maxAge": 30,
    "maxBackups": 5,
    "compress": true
  },
  "discord": {
    "channelId": "https://discord.com/api/webhooks/YOUR_WEBHOOK_URL"
  },
  "commands": [
    {
      "name": "Show Docker Status",
      "description": "Lists all running Docker containers",
      "command": "docker",
      "args": ["ps", "-a"]
    },
    {
      "name": "Git Status",
      "description": "Shows the working tree status",
      "command": "git",
      "args": ["status"]
    }
  ]
}
```

### Configuration Options

#### Main Configuration

| Field | Description | Default | Required |
|-------|-------------|---------|----------|
| `workingDir` | Global working directory for commands | Current directory | No |
| `docker.host` | Docker daemon socket | `unix:///var/run/docker.sock` | No |
| `discord.channelId` | Discord webhook URL | None | Yes |
| `commands` | Array of commands to execute | [] | Yes |

#### Logging Configuration (Optional)

| Field | Description | Default |
|-------|-------------|--------|
| `logs.directory` | Directory to store log files | `./logs` |
| `logs.maxSize` | Maximum size of log files in MB | 10 |
| `logs.maxAge` | Maximum age of log files in days | 30 |
| `logs.maxBackups` | Maximum number of old log files to keep | 5 |
| `logs.compress` | Whether to compress old log files | true |

#### Command Structure

| Field | Description | Required |
|-------|-------------|----------|
| `name` | Name of the command | Yes |
| `description` | Description of what the command does | Yes |
| `command` | The executable to run | Yes |
| `args` | Array of arguments to pass to the command | No |
| `dir` | Working directory specific to this command | No |
| `envVars` | Environment variables for the command | No |

### Discord Integration

Delivr works with Discord webhooks. Simply create a webhook in your Discord channel and paste the URL in the `channelId` field of your configuration file.

To create a webhook in Discord:

1. Go to your Discord server
2. Edit a channel
3. Select 'Integrations'
4. Click 'Webhooks'
5. Click 'New Webhook'
6. Copy the webhook URL

## Environment Variables

- `DELIVR_CONFIG`: Path to the config file (overrides the default location)

## GitHub Actions Integration

This project includes a GitHub Actions workflow that automatically builds the application for Linux (AMD64) on each push or pull request to the main branch. Tagged versions will include the tag name in the built binary filename.

The workflow is defined in `.github/workflows/build.yml`.

## Log Files

Each command generates its own log file in the format `command-name-YYYY-MM-DD.log`. These logs contain:

- Command execution details (time, arguments, working directory)
- Complete stdout and stderr output
- Execution status and duration

## Minimal Configuration Example

The following is a minimal configuration example with only the required fields:

```yaml
discord:
  channelId: https://discord.com/api/webhooks/YOUR_WEBHOOK_URL

commands:
  - name: Git Status
    description: Shows git status
    command: git
    args: [status]
```

## Credits

This project was developed with the assistance of:

- [Windsurf](https://windsurf.com) - Agentic AI coding platform
- [Claude](https://www.anthropic.com/claude) - AI assistant by Anthropic

## License

MIT
