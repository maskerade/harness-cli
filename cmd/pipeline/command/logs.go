package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/harness/harness-cli/cmd/cmdutils"
	"github.com/harness/harness-cli/config"

	"github.com/spf13/cobra"
)

// NewGetLogsCmd creates the logs command
func NewGetLogsCmd(f *cmdutils.Factory) *cobra.Command {
	var outputDir string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "logs EXECUTION_ID",
		Short: "Download execution logs",
		Long:  "Downloads logs for a specific pipeline execution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			executionID := args[0]

			// Get required parameters
			orgID := config.Global.OrgID
			projectID := config.Global.ProjectID

			if orgID == "" {
				return fmt.Errorf("organization ID is required (use --org flag or auth login)")
			}
			if projectID == "" {
				return fmt.Errorf("project ID is required (use --project flag or auth login)")
			}

			// Create pipeline client
			client := f.PipelineClient()

			// Get logs
			logs, err := client.GetExecutionLogs(context.Background(), orgID, projectID, executionID)
			if err != nil {
				return fmt.Errorf("failed to get logs: %w", err)
			}

			// Determine output location
			if outputFile == "" {
				outputFile = fmt.Sprintf("%s-logs.txt", executionID)
			}

			if outputDir != "" {
				// Create output directory if it doesn't exist
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf("failed to create output directory: %w", err)
				}
				outputFile = filepath.Join(outputDir, outputFile)
			}

			// Write logs to file
			if err := os.WriteFile(outputFile, logs, 0644); err != nil {
				return fmt.Errorf("failed to write logs: %w", err)
			}

			fmt.Printf("Logs saved to: %s\n", outputFile)
			fmt.Printf("Log size: %d bytes\n", len(logs))

			return nil
		},
	}

	cmd.Flags().StringVar(&outputDir, "output-dir", "", "directory to save logs (default: current directory)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file name (default: EXECUTION_ID-logs.txt)")

	return cmd
}
