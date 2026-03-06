# Pipeline Commands for Harness CLI

## Overview

This extension adds comprehensive pipeline management capabilities to the Harness CLI tool. These commands enable you to interact with Harness CI/CD pipelines, check pipeline execution status, download logs, trigger pipelines, and manage pipeline configurations.

## Installation

The pipeline commands are integrated into the main Harness CLI. Build the CLI from source:

```bash
git clone https://github.com/harness/harness-cli.git
cd harness-cli
go build -o hc ./cmd/hc
```

## Authentication

Before using pipeline commands, authenticate with Harness:

```bash
hc auth login
```

Provide your:
- API URL (default: https://app.harness.io)
- API Token (Personal Access Token from Harness)
- Organization ID
- Project ID

## Command Reference

### List Pipelines

Lists all pipelines in the specified organization and project.

```bash
hc pipeline list [flags]
hc pipe list     # using alias
```

**Flags:**
- `--page-size` - Number of items per page (default: 20)
- `--page` - Page number, zero-indexed (default: 0)
- `--org` - Organization ID (override saved config)
- `--project` - Project ID (override saved config)
- `--format` - Output format: `table` or `json` (default: table)

**Example:**
```bash
hc pipeline list
hc pipeline list --page-size 50 --format json
hc pipeline list --org myorg --project myproj
```

**Output Columns:**
- Pipeline ID
- Name
- Stages (count)
- Last Status
- Last Updated

---

### Get Pipeline Details

Retrieves detailed information about a specific pipeline.

```bash
hc pipeline get PIPELINE_ID [flags]
```

**Flags:**
- `--yaml` - Show pipeline YAML definition
- `--org` - Organization ID
- `--project` - Project ID
- `--format` - Output format: `table` or `json`

**Example:**
```bash
hc pipeline get my-pipeline
hc pipeline get my-pipeline --yaml
hc pipeline get my-pipeline --format json
```

**Output:**
- Pipeline name and identifier
- Description
- Created and updated timestamps
- Tags
- YAML definition (with --yaml flag)

---

### List Executions

Lists pipeline executions with optional filters.

```bash
hc pipeline executions [PIPELINE_ID] [flags]
```

**Arguments:**
- `PIPELINE_ID` - Optional. Filter by specific pipeline

**Flags:**
- `--status` - Filter by execution status: `SUCCESS`, `FAILED`, `RUNNING`, `ABORTED`
- `--page-size` - Number of items per page (default: 20)
- `--page` - Page number (default: 0)
- `--org` - Organization ID
- `--project` - Project ID
- `--format` - Output format

**Examples:**
```bash
# List all executions
hc pipeline executions

# List executions for specific pipeline
hc pipeline executions my-pipeline

# List only failed executions
hc pipeline executions --status FAILED

# List running executions for a pipeline
hc pipeline executions my-pipeline --status RUNNING

# Get more results
hc pipeline executions --page-size 50
```

**Output Columns:**
- Execution ID
- Pipeline
- Name
- Status
- Started (timestamp)
- Ended (timestamp)

---

### Get Execution Details

Retrieves detailed information about a specific pipeline execution.

```bash
hc pipeline execution EXECUTION_ID [flags]
```

**Flags:**
- `--org` - Organization ID
- `--project` - Project ID
- `--format` - Output format

**Example:**
```bash
hc pipeline execution abc123xyz
hc pipeline execution abc123xyz --format json
```

**Output:**
- Execution ID
- Pipeline name and ID
- Status
- Start and end timestamps
- Duration
- Trigger information
- Stage names

---

### Download Execution Logs

Downloads logs for a specific pipeline execution.

```bash
hc pipeline logs EXECUTION_ID [flags]
```

**Flags:**
- `--output-dir` - Directory to save logs (default: current directory)
- `--output`, `-o` - Output file name (default: EXECUTION_ID-logs.txt)
- `--org` - Organization ID
- `--project` - Project ID

**Examples:**
```bash
# Save to current directory
hc pipeline logs abc123xyz

# Save to specific directory
hc pipeline logs abc123xyz --output-dir ./logs

# Custom filename
hc pipeline logs abc123xyz --output build-logs.txt

# Both directory and filename
hc pipeline logs abc123xyz --output-dir ./logs --output build.log
```

**Output:**
- Log file saved confirmation
- File path
- Log size in bytes

---

### Trigger Pipeline

Triggers a pipeline execution with optional input sets and runtime inputs.

**⚠️ Important:** Requires `--confirm` flag for safety.

```bash
hc pipeline trigger PIPELINE_ID [flags]
```

**Flags:**
- `--confirm` - Required. Confirm pipeline trigger (prevents accidental executions)
- `--input-set` - Input set references (can be specified multiple times)
- `--runtime-input` - Path to runtime input YAML file
- `--org` - Organization ID
- `--project` - Project ID

**Examples:**
```bash
# Trigger pipeline (dry-run preview)
hc pipeline trigger my-pipeline

# Confirm and trigger
hc pipeline trigger my-pipeline --confirm

# Trigger with input set
hc pipeline trigger my-pipeline --input-set prod-inputs --confirm

# Trigger with multiple input sets
hc pipeline trigger my-pipeline \
  --input-set common-vars \
  --input-set prod-config \
  --confirm

# Trigger with runtime inputs from file
hc pipeline trigger my-pipeline \
  --runtime-input runtime-vars.yaml \
  --confirm
```

**Output:**
- Confirmation message
- Execution ID
- Instructions to check status

---

### List Input Sets

Lists all input sets for a specific pipeline.

```bash
hc pipeline input-sets PIPELINE_ID [flags]
```

**Flags:**
- `--org` - Organization ID
- `--project` - Project ID
- `--format` - Output format

**Example:**
```bash
hc pipeline input-sets my-pipeline
hc pipeline input-sets my-pipeline --format json
```

**Output Columns:**
- Input Set ID
- Name
- Type
- Description
- Last Updated

---

### List Triggers

Lists all triggers configured for a specific pipeline.

```bash
hc pipeline triggers PIPELINE_ID [flags]
```

**Flags:**
- `--org` - Organization ID
- `--project` - Project ID
- `--format` - Output format

**Example:**
```bash
hc pipeline triggers my-pipeline
hc pipeline triggers my-pipeline --format json
```

**Output Columns:**
- Trigger ID
- Name
- Type
- Enabled (true/false)
- Description

---

## Common Workflows

### Debug a Failed Pipeline

```bash
# 1. List recent failed executions
hc pipeline executions --status FAILED --page-size 10

# 2. Get details of the failed execution
hc pipeline execution <EXECUTION_ID>

# 3. Download logs for analysis
hc pipeline logs <EXECUTION_ID> --output-dir ./debug-logs

# 4. View pipeline configuration
hc pipeline get <PIPELINE_ID> --yaml
```

### Monitor Pipeline Status

```bash
# Check all running pipelines
hc pipeline executions --status RUNNING

# Watch specific pipeline executions
watch -n 5 'hc pipeline executions my-pipeline --status RUNNING'

# Get execution details
hc pipeline execution <EXECUTION_ID>
```

### Trigger Pipeline with Configuration

```bash
# 1. List available input sets
hc pipeline input-sets my-pipeline

# 2. Review pipeline configuration
hc pipeline get my-pipeline --yaml

# 3. Trigger with selected input set
hc pipeline trigger my-pipeline \
  --input-set prod-config \
  --confirm

# 4. Monitor execution
hc pipeline execution <EXECUTION_ID>
```

### Audit Pipeline Triggers

```bash
# List all triggers for a pipeline
hc pipeline triggers my-pipeline

# Check recent executions
hc pipeline executions my-pipeline --page-size 20

# Review execution details to see trigger info
hc pipeline execution <EXECUTION_ID> --format json
```

---

## Global Flags

Available for all pipeline commands:

- `--account` - Account ID (overrides saved config)
- `--org` - Organization ID (overrides saved config)
- `--project` - Project ID (overrides saved config)
- `--token` - Authentication token (overrides saved config)
- `--api-url` - Base URL for the API (overrides saved config)
- `--format` - Output format: `table` or `json` (default: table)
- `--verbose`, `-v` - Enable verbose logging

---

## Environment Variables

You can set these instead of using flags:

```bash
export HARNESS_API_KEY="your-api-token"
export HARNESS_API_URL="https://app.harness.io"
export HARNESS_ORG_ID="your-org-id"
export HARNESS_PROJECT_ID="your-project-id"
```

---

## Error Handling

### Common Errors

**"Not logged in. Please run 'hc auth login' first"**
- Solution: Run `hc auth login` to authenticate

**"organization ID is required"**
- Solution: Provide `--org` flag or set during login

**"project ID is required"**
- Solution: Provide `--project` flag or set during login

**"API error (status 401)"**
- Solution: Re-authenticate with `hc auth login`

**"API error (status 403)"**
- Solution: Check token permissions in Harness

**"API error (status 404)"**
- Solution: Verify pipeline/execution ID exists

---

## Output Formats

### Table Format (Default)

```
Pipeline ID          Name              Stages  Last Status  Last Updated
my-pipeline          Production Build  5       SUCCESS      1709740800
api-deploy           API Deployment    3       FAILED       1709737200
```

### JSON Format

```bash
hc pipeline list --format json
```

```json
{
  "content": [
    {
      "identifier": "my-pipeline",
      "name": "Production Build",
      "numOfStages": 5,
      "executionSummaryInfo": {
        "lastExecutionStatus": "SUCCESS"
      },
      "lastUpdatedAt": 1709740800
    }
  ]
}
```

---

## Tips & Best Practices

1. **Use aliases**: `hc pipe` is shorter than `hc pipeline`

2. **Save credentials**: Use `hc auth login` to save credentials instead of passing flags

3. **JSON for scripting**: Use `--format json` when using in scripts or CI/CD

4. **Confirm safety**: `trigger` command requires `--confirm` to prevent accidents

5. **Page through results**: Use `--page-size` and `--page` for large result sets

6. **Monitor with watch**: Use `watch` command to monitor running pipelines

7. **Download logs early**: Download logs before they're rotated out

8. **Use input sets**: Manage different configurations as input sets

---

## API Reference

The pipeline commands use the Harness Pipeline API:
- Base path: `/pipeline/api/`
- Authentication: x-api-key header
- Format: JSON

### Endpoints Used

- `POST /pipeline/api/pipelines/list` - List pipelines
- `GET /pipeline/api/pipelines/{id}` - Get pipeline
- `POST /pipeline/api/pipelines/execution/summary` - List executions
- `GET /pipeline/api/pipelines/execution/{id}` - Get execution
- `GET /pipeline/api/pipelines/execution/{id}/logs` - Download logs
- `POST /pipeline/api/pipelines/execution/{id}` - Trigger pipeline
- `GET /pipeline/api/inputSets` - List input sets
- `GET /pipeline/api/triggers` - List triggers

---

## Contributing

Found a bug or want to add a feature? Contributions are welcome!

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

---

## License

MIT License - See LICENSE file for details

---

## Support

- **Issues:** [GitHub Issues](https://github.com/harness/harness-cli/issues)
- **Documentation:** [Harness Docs](https://docs.harness.io)
- **Community:** [Harness Community](https://community.harness.io)
