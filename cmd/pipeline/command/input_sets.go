package command

import (
	"context"
	"fmt"

	"github.com/harness/harness-cli/cmd/cmdutils"
	"github.com/harness/harness-cli/config"
	"github.com/harness/harness-cli/util/common/printer"

	"github.com/spf13/cobra"
)

// NewListInputSetsCmd creates the list input sets command
func NewListInputSetsCmd(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "input-sets PIPELINE_ID",
		Short: "List pipeline input sets",
		Long:  "Lists all input sets for a specific pipeline",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pipelineID := args[0]

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

			// List input sets
			response, err := client.ListInputSets(context.Background(), orgID, projectID, pipelineID)
			if err != nil {
				return fmt.Errorf("failed to list input sets: %w", err)
			}

			if response.Status != "SUCCESS" {
				return fmt.Errorf("API returned non-success status: %s", response.Status)
			}

			// Print results
			totalItems := int64(len(response.Data.Content))

			err = printer.Print(response.Data.Content, 0,
				1, totalItems, totalItems > 0, [][]string{
					{"identifier", "Input Set ID"},
					{"name", "Name"},
					{"inputSetType", "Type"},
					{"description", "Description"},
					{"lastUpdatedAt", "Last Updated"},
				})

			return err
		},
	}

	return cmd
}
