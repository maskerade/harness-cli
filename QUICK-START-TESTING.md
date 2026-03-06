# Quick Start: Testing Pipeline Commands

## Prerequisites

- Go 1.20+ installed
- Access to a Harness account
- Harness Personal Access Token (PAT)

## Step 1: Build the CLI

```bash
cd harness-cli
go mod tidy
go build -o hc ./cmd/hc
```

## Step 2: Authenticate

```bash
./hc auth login
```

You'll be prompted for:
- **API URL:** https://app.harness.io (or your self-hosted URL)
- **API Token:** Your Harness PAT
- **Organization ID:** Your org identifier
- **Project ID:** Your project identifier

## Step 3: Verify Authentication

```bash
./hc auth status
```

Should show your logged-in status and configuration.

## Step 4: Test Pipeline Commands

### List Pipelines

```bash
./hc pipeline list
```

**Expected:** Table showing your pipelines with columns:
- Pipeline ID
- Name
- Stages
- Last Status
- Last Updated

### Get Pipeline Details

```bash
# Replace with your actual pipeline ID
./hc pipeline get your-pipeline-id
```

**Expected:** Detailed pipeline information including description, tags, timestamps

### List Executions

```bash
# All executions
./hc pipeline executions

# Specific pipeline
./hc pipeline executions your-pipeline-id

# Only failed executions
./hc pipeline executions --status FAILED

# Only running executions
./hc pipeline executions --status RUNNING
```

**Expected:** Table showing executions with ID, pipeline, name, status, timestamps

### Get Execution Details

```bash
# Replace with actual execution ID from previous command
./hc pipeline execution your-execution-id
```

**Expected:** Detailed execution info including status, duration, trigger info, stages

### Download Logs

```bash
# Download logs for an execution
./hc pipeline logs your-execution-id

# Save to specific directory
./hc pipeline logs your-execution-id --output-dir ./logs

# Custom filename
./hc pipeline logs your-execution-id --output my-logs.txt
```

**Expected:**
- Log file saved message
- File path confirmation
- Log size in bytes

### List Input Sets

```bash
./hc pipeline input-sets your-pipeline-id
```

**Expected:** Table of input sets with ID, name, type, description

### List Triggers

```bash
./hc pipeline triggers your-pipeline-id
```

**Expected:** Table of triggers with ID, name, type, enabled status

### Trigger Pipeline (With Confirmation)

```bash
# Dry run (shows what would be triggered)
./hc pipeline trigger your-pipeline-id

# Actual trigger (requires --confirm)
./hc pipeline trigger your-pipeline-id --confirm

# With input set
./hc pipeline trigger your-pipeline-id --input-set prod-config --confirm
```

**Expected:**
- Execution ID returned
- Success confirmation
- Instructions to check status

## Step 5: Test Output Formats

### JSON Format

```bash
./hc pipeline list --format json
./hc pipeline get your-pipeline-id --format json
./hc pipeline executions --format json
```

**Expected:** Pretty-printed JSON output

### Table Format (Default)

```bash
./hc pipeline list --format table
```

**Expected:** Formatted table with columns

## Step 6: Test Aliases

```bash
# These should work the same
./hc pipeline list
./hc pipe list
```

## Step 7: Test Help

```bash
# Main pipeline help
./hc pipeline --help
./hc pipe --help

# Individual command help
./hc pipeline list --help
./hc pipeline get --help
./hc pipeline executions --help
./hc pipeline execution --help
./hc pipeline logs --help
./hc pipeline trigger --help
./hc pipeline input-sets --help
./hc pipeline triggers --help
```

**Expected:** Usage information and flag descriptions for each command

## Step 8: Test Global Flags

```bash
# Override org/project
./hc pipeline list --org another-org --project another-project

# Use different account
./hc pipeline list --account another-account

# Verbose output
./hc pipeline list --verbose

# Different API URL
./hc pipeline list --api-url https://custom.harness.io
```

## Step 9: Test Error Handling

### Without Authentication

```bash
./hc auth logout
./hc pipeline list
```

**Expected:** Error message: "Not logged in. Please run 'hc auth login' first"

### With Missing Parameters

```bash
./hc auth login --non-interactive --token YOUR_TOKEN
# (don't provide org/project)

./hc pipeline list
```

**Expected:** Error message about missing org/project ID

### Invalid Pipeline ID

```bash
./hc pipeline get nonexistent-pipeline
```

**Expected:** API error 404 message

## Common Test Scenarios

### Scenario 1: Monitor Running Pipeline

```bash
# Start monitoring
watch -n 5 './hc pipeline executions --status RUNNING'

# In another terminal, check specific execution
./hc pipeline execution <execution-id>

# Download logs when complete
./hc pipeline logs <execution-id>
```

### Scenario 2: Debug Failed Pipeline

```bash
# Find recent failures
./hc pipeline executions --status FAILED --page-size 5

# Get failure details
./hc pipeline execution <failed-execution-id> --format json

# Download failure logs
./hc pipeline logs <failed-execution-id> --output-dir ./failure-logs

# Check pipeline configuration
./hc pipeline get <pipeline-id> --yaml
```

### Scenario 3: Automated Pipeline Trigger

```bash
# Create script: trigger-and-wait.sh
#!/bin/bash
set -e

PIPELINE_ID="$1"
EXECUTION_ID=$(./hc pipeline trigger "$PIPELINE_ID" --confirm --format json | jq -r '.data.planExecutionId')

echo "Triggered execution: $EXECUTION_ID"
echo "Waiting for completion..."

while true; do
  STATUS=$(./hc pipeline execution "$EXECUTION_ID" --format json | jq -r '.data.pipelineExecutionSummary.status')
  echo "Status: $STATUS"

  if [[ "$STATUS" != "RUNNING" ]]; then
    break
  fi

  sleep 10
done

echo "Execution $STATUS"
./hc pipeline logs "$EXECUTION_ID" --output-dir ./logs
```

## Troubleshooting

### Command Not Found

```bash
# Make sure binary is executable
chmod +x hc

# Or use full path
/path/to/hc pipeline list
```

### API Errors

```bash
# Check auth status
./hc auth status

# Re-authenticate
./hc auth logout
./hc auth login

# Verbose mode for debugging
./hc pipeline list --verbose
```

### Permission Errors

Ensure your Harness PAT has:
- `core_pipeline_view` - Read pipeline definitions
- `core_execution_view` - Read pipeline executions
- `core_pipeline_execute` - Trigger pipelines (for trigger command)

### Build Errors

```bash
# Clean and rebuild
go clean
go mod tidy
go build -o hc ./cmd/hc
```

## Success Indicators

✅ All commands show help correctly
✅ List commands return data (or empty list if no pipelines)
✅ Get commands return 404 for invalid IDs
✅ Format flag works (table vs JSON)
✅ Aliases work (pipe = pipeline)
✅ Trigger requires --confirm flag
✅ Logs download successfully
✅ Error messages are clear and helpful

## Next Steps

After successful testing:

1. **Create Test Automation**
   - Add unit tests for client
   - Add integration tests for commands
   - Add CI/CD pipeline for testing

2. **Gather Feedback**
   - Share with team
   - Test with real workflows
   - Document common use cases

3. **Iterate**
   - Address issues found
   - Add requested features
   - Improve error messages

## Resources

- **User Guide:** PIPELINE-COMMANDS.md
- **Implementation Summary:** IMPLEMENTATION-SUMMARY.md
- **Harness API Docs:** https://apidocs.harness.io/
- **CLI Source:** https://github.com/harness/harness-cli

## Support

If you encounter issues:

1. Check verbose output: `--verbose` flag
2. Verify authentication: `hc auth status`
3. Check API permissions in Harness
4. Review error messages carefully
5. Open GitHub issue with details

---

**Happy Testing! 🚀**
