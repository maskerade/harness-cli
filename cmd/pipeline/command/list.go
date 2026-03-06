package command

import (
	"context"
	"fmt"

	"github.com/harness/harness-cli/cmd/cmdutils"
	"github.com/harness/harness-cli/config"
	"github.com/harness/harness-cli/util/common/printer"

	"github.com/spf13/cobra"
)

// NewListPipelineCmd creates the list command for pipelines
func NewListPipelineCmd(f *cmdutils.Factory) *cobra.Command {
	var pageSize int
	var pageIndex int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all pipelines",
		Long:  "Lists all pipelines in the specified organization and project",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			// List pipelines
			response, err := client.ListPipelines(context.Background(), orgID, projectID, pageSize, pageIndex)
			if err != nil {
				return fmt.Errorf("failed to list pipelines: %w", err)
			}

			if response.Status != "SUCCESS" {
				return fmt.Errorf("API returned non-success status: %s", response.Status)
			}

			// Print results using the printer utility
			totalPages := int64(response.Data.TotalPages)
			totalItems := int64(response.Data.TotalItems)
			pageItems := response.Data.PageItemCount

			err = printer.Print(response.Data.Content, int64(response.Data.PageIndex),
				totalPages, totalItems, pageItems > 0, [][]string{
					{"identifier", "Pipeline ID"},
					{"name", "Name"},
					{"numOfStages", "Stages"},
					{"executionSummaryInfo.lastExecutionStatus", "Last Status"},
					{"lastUpdatedAt", "Last Updated"},
				})

			return err
		},
	}

	cmd.Flags().IntVar(&pageSize, "page-size", 20, "number of items per page")
	cmd.Flags().IntVar(&pageIndex, "page", 0, "page number (zero-indexed)")

	return cmd
}
