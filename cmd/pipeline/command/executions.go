package command

import (
	"context"
	"fmt"

	"github.com/harness/harness-cli/cmd/cmdutils"
	"github.com/harness/harness-cli/config"
	"github.com/harness/harness-cli/util/common/printer"

	"github.com/spf13/cobra"
)

// NewListExecutionsCmd creates the list executions command
func NewListExecutionsCmd(f *cmdutils.Factory) *cobra.Command {
	var pipelineID string
	var status string
	var pageSize int
	var pageIndex int

	cmd := &cobra.Command{
		Use:   "executions [PIPELINE_ID]",
		Short: "List pipeline executions",
		Long:  "Lists pipeline executions. Optionally filter by pipeline ID and status",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get pipeline ID from args if provided
			if len(args) > 0 {
				pipelineID = args[0]
			}

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

			// List executions
			response, err := client.ListExecutions(context.Background(), orgID, projectID, pipelineID, status, pageSize, pageIndex)
			if err != nil {
				return fmt.Errorf("failed to list executions: %w", err)
			}

			if response.Status != "SUCCESS" {
				return fmt.Errorf("API returned non-success status: %s", response.Status)
			}

			// Print results
			totalPages := int64(response.Data.TotalPages)
			totalItems := int64(response.Data.TotalItems)
			pageItems := response.Data.PageItemCount

			err = printer.Print(response.Data.Content, int64(response.Data.PageIndex),
				totalPages, totalItems, pageItems > 0, [][]string{
					{"planExecutionId", "Execution ID"},
					{"pipelineIdentifier", "Pipeline"},
					{"name", "Name"},
					{"status", "Status"},
					{"startTs", "Started"},
					{"endTs", "Ended"},
				})

			return err
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "filter by execution status (SUCCESS, FAILED, RUNNING, ABORTED)")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "number of items per page")
	cmd.Flags().IntVar(&pageIndex, "page", 0, "page number (zero-indexed)")

	return cmd
}
