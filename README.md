# Harness CLI Tool (hc)

A powerful command-line interface tool for interacting with Harness services

## Overview

The Harness CLI (hc) provides a unified command-line interface for interacting with various Harness services. It follows a consistent, resource-based command structure:

```
hc [<global-flags>] <command> <subcommand> [<positional-args>…] [<flags>]
```

### Available Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `auth` | - | Authentication commands (login, logout, status) |
| `registry` | `reg` | Manage Harness Artifact Registries |
| `artifact` | `art` | Manage artifacts in registries |
| `pipeline` | `pipe` | Manage Harness CI/CD Pipelines |
| `project` | `proj` | Manage Harness Projects |
| `organisation` | `org` | Manage Harness Organisations |
| `api` | - | Raw REST API passthrough for power users |

## Installation

### Quick Install (Recommended)

Install the latest version with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/harness/harness-cli/v2/install | sh
```

Or with sudo if you need elevated privileges:

```bash
curl -fsSL https://raw.githubusercontent.com/harness/harness-cli/v2/install | sudo sh
```

This script automatically detects your OS and architecture, downloads the appropriate binary, verifies its checksum for security, and installs it to `/usr/local/bin`.

### Custom Installation Directory

You can install to a custom directory by setting the `INSTALL_DIR` environment variable:

```bash
curl -fsSL https://raw.githubusercontent.com/harness/harness-cli/v2/install | INSTALL_DIR=$HOME/.local/bin sh
```

### Manual Binary Installation

Download the latest binary from the [releases page](https://github.com/harness/harness-cli/releases):

```bash
# Download the latest release for your platform
# Make it executable
chmod +x hc
# Move it to a directory in your PATH
mv hc /usr/local/bin/
```

### Building from Source

```bash
# Install go if you haven't

# Clone the repository
git clone https://github.com/harness/harness-cli.git
cd harness-cli

# Build the binary
make build
```

## Configuration

The CLI can be configured using:

1. Configuration file at `$HOME/.harness/auth.json`
2. Environment variables (coming soon)
3. Command-line flags

### Authentication

Before using most commands, you need to authenticate:

```bash
# Login with API key
hc auth login

# Check authentication status
hc auth status

# Logout
hc auth logout
```

You can also provide credentials via:
- Configuration file at `$HOME/.harness/auth.json`
- Command-line flags: `--token`, `--account`, `--org`, `--project`
- Environment variables

## Command Reference

### Authentication (`hc auth`)

Manage authentication with Harness services.

```bash
# Login interactively
hc auth login

# Login with API key
hc auth login --api-key <your-api-key>

# Check authentication status
hc auth status

# Logout
hc auth logout
```

### Registry Management (`hc registry` or `hc reg`)

Manage Harness Artifact Registries.

```bash
# List all registries
hc registry list
hc reg list  # Using alias

# Get registry details
hc registry get <registry-name>

# Create a registry (coming soon)
hc registry create <registry-name> --package-type DOCKER

# Delete a registry
hc registry delete <registry-name>

# Migrate artifacts from external registries
hc registry migrate --config migrate-config.yaml
```

### Artifact Management (`hc artifact` or `hc art`)

Manage artifacts within registries.

```bash
# List all artifacts
hc artifact list
hc art list  # Using alias

# List artifacts in a specific registry
hc artifact list --registry <registry-name>

# Delete an artifact (deletes all versions)
hc artifact delete <artifact-name> --registry <registry-name>

# Delete a specific version of an artifact
hc artifact delete <artifact-name> --registry <registry-name> --version <version>

# Push artifacts
hc artifact push generic <registry-name> <file-path> --name <artifact-name> --version <version>
hc artifact push go <registry-name> <module-path>

# Pull artifacts
hc artifact pull generic <registry-name> <package-path> <destination>
```

### Pipeline Management (`hc pipeline` or `hc pipe`)

Manage Harness CI/CD Pipelines, monitor executions, and trigger pipeline runs.

```bash
# List all pipelines
hc pipeline list
hc pipe list  # Using alias

# List with pagination
hc pipeline list --page 2 --size 50

# Get pipeline details
hc pipeline get <pipeline-id>

# Get pipeline YAML definition
hc pipeline get <pipeline-id> --yaml

# List pipeline executions
hc pipeline executions

# List executions for a specific pipeline
hc pipeline executions --pipeline <pipeline-id>

# Filter executions by status
hc pipeline executions --status Success
hc pipeline executions --status Failed --pipeline <pipeline-id>

# Get execution details
hc pipeline execution <execution-id>

# Download execution logs
hc pipeline logs <execution-id>

# Save logs to a custom file
hc pipeline logs <execution-id> --output execution-logs.txt
hc pipeline logs <execution-id> -o logs.txt

# Save logs to a directory
hc pipeline logs <execution-id> --output-dir ./logs

# Trigger a pipeline (requires --confirm flag for safety)
hc pipeline trigger <pipeline-id> --confirm

# Trigger with input sets
hc pipeline trigger <pipeline-id> --input-set prod-config --confirm
hc pipeline trigger <pipeline-id> --input-set common --input-set prod --confirm

# Trigger with runtime input YAML
hc pipeline trigger <pipeline-id> --runtime-input inputs.yaml --confirm

# List pipeline input sets
hc pipeline input-sets <pipeline-id>

# List pipeline triggers
hc pipeline triggers <pipeline-id>
```

#### Pipeline Command Details

**List Pipelines** (`hc pipeline list`)
- Lists all pipelines in your project
- Supports pagination with `--page` and `--size` flags
- Default page size is 20

**Get Pipeline** (`hc pipeline get`)
- Get detailed information about a specific pipeline
- Use `--yaml` flag to output the pipeline YAML definition
- Useful for inspecting pipeline configuration

**List Executions** (`hc pipeline executions`)
- Lists pipeline execution history
- Filter by specific pipeline using `--pipeline` flag
- Filter by status: `Success`, `Failed`, `Running`, `Aborted`, etc.
- Supports pagination for large result sets

**Get Execution** (`hc pipeline execution`)
- Get detailed information about a specific execution
- Shows execution status, timing, and stage information
- Use execution ID from the executions list

**Download Logs** (`hc pipeline logs`)
- Downloads complete execution logs
- Default saves to `execution-<id>-logs.txt`
- Use `--output` or `-o` for custom filename
- Use `--output-dir` to specify save directory

**Trigger Pipeline** (`hc pipeline trigger`)
- Manually trigger a pipeline execution
- **Requires `--confirm` flag** as a safety measure
- Optionally specify input sets with `--input-set` (can be used multiple times)
- Provide runtime inputs via `--runtime-input` YAML file
- Returns execution ID for tracking

**List Input Sets** (`hc pipeline input-sets`)
- Lists all input sets configured for a pipeline
- Shows input set identifiers, names, and types

**List Triggers** (`hc pipeline triggers`)
- Lists all triggers configured for a pipeline
- Shows trigger names, types (Webhook, Scheduled, etc.), and enabled status

### Project Management (`hc project` or `hc proj`) (coming soon)

Manage Harness Projects.

```bash
# List all projects
hc project list

# Get project details
hc project get <project-id>

# Create a project (coming soon)
hc project create <project-id>

# Delete a project (coming soon)
hc project delete <project-id>
```

### Organisation Management (`hc organisation` or `hc org`) (coming soon)

Manage Harness Organisations.

```bash
# List all organisations
hc organisation list
hc org list  # Using alias

# Get organisation details
hc org get <org-id>

# Create an organisation (coming soon)
hc org create <org-id>

# Delete an organisation (coming soon)
hc org delete <org-id>
```

### API Passthrough (`hc api`) (coming soon)

Make raw REST API calls to Harness (for power users).

```bash
# GET request
hc api /har/api/v1/registries

# POST request with data
hc api /har/api/v1/registries --method POST --data '{"identifier":"my-registry"}'

# Custom headers
hc api /har/api/v1/registries --header "Content-Type: application/json"

# PUT/DELETE requests
hc api /har/api/v1/registries/my-reg --method DELETE
```

## Global Flags

The following flags are available for all commands:

```bash
--account string          Account ID (overrides saved config)
--api-url string          Base URL for the API (overrides saved config)
--token string            Authentication token (overrides saved config)
--org string              Organisation ID (overrides saved config)
--project string          Project ID (overrides saved config)
--format string           Output format: table (default) or json
--log-file string         Path to store logs
```

## Output Formatting 

The CLI supports different output formats using the `--format` flag:

```bash
# Output in JSON format
hc registry list --format=json

# Output in table format (default)
hc registry list --format=table

# Works with all list/get commands
hc artifact list --registry my-reg --format=json
```

JSON output supports:
- Pretty printing with configurable indentation
- Smart pagination information
- Custom output formatting

## Development

### Project Structure

```
harness-cli/
├── api/              # OpenAPI specs for each service
├── cmd/              # CLI commands implementation
│   ├── hc/           # Main CLI entry point
│   ├── auth/         # Authentication commands
│   ├── registry/     # Registry management commands
│   ├── artifact/     # Artifact management commands
│   ├── pipeline/     # Pipeline management commands
│   ├── project/      # Project management commands
│   ├── organisation/ # Organisation management commands
│   └── api/          # API passthrough command
├── config/           # Configuration handling
├── internal/         # Internal packages and generated API clients
├── module/           # Service-specific modules
├── tools/            # Development tools
└── util/             # Utility functions
```

### Adding New Commands (coming soon)
TODO

### Building

```bash
# Build the binary
make build

# Run tests
make test

# Run linter
make lint
```

## License

MIT License
