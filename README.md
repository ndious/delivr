# Delivr

A simple Go application that executes commands (Docker, Git, or any other CLI tool) and sends notifications to Discord.

> This project was built with [Windsurf](https://windsurf.com) and [Claude](https://www.anthropic.com/claude) AI assistance.

## Features

- Execute various types of commands (Docker, Git, etc.)
- Send command output to Discord via bot token or webhook
- Log command outputs to files with automatic rotation
- Configure via JSON with support for custom paths
- Error handling and comprehensive reporting
- One-time execution or daemon mode

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

Edit the generated `config.json` file with your Discord token or webhook URL and customize your commands.

## Usage

```bash
# Run with default config file (config.json in current directory or ~/.delivr/config.json)
./delivr

# Specify a custom config file
./delivr --config /path/to/your/config.json

# Run in daemon mode (doesn't exit after executing all commands)
./delivr --daemon

# Generate a default configuration file
./delivr --init

# Generate a configuration file at a specific location
./delivr --init --out /path/to/new/config.json
```

## Configuration

The configuration file has the following structure:

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
    "token": "YOUR_DISCORD_BOT_TOKEN_HERE",
    "channelId": "YOUR_DISCORD_CHANNEL_ID_OR_WEBHOOK_URL"
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

| Field | Description | Default |
|-------|-------------|--------|
| `workingDir` | Global working directory for commands | Current directory |
| `docker.host` | Docker daemon socket | `unix:///var/run/docker.sock` |
| `commands` | Array of commands to execute | [] |

#### Logging Configuration

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

Delivr supports two methods of Discord integration:

1. **Discord Webhook URL**: The easiest method. Simply create a webhook in your Discord channel and paste the URL in the `channelId` field.

2. **Discord Bot**: For more advanced functionality (not currently implemented). Create a bot through the Discord Developer Portal and paste the token in the `token` field.

## Environment Variables

- `DELIVR_CONFIG`: Path to the config file (overrides the default location)

## Log Files

Each command generates its own log file in the format `command-name-YYYY-MM-DD.log`. These logs contain:

- Command execution details (time, arguments, working directory)
- Complete stdout and stderr output
- Execution status and duration

## Credits

This project was developed with the assistance of:

- [Windsurf](https://windsurf.com) - Agentic AI coding platform
- [Claude](https://www.anthropic.com/claude) - AI assistant by Anthropic

## License

MIT
