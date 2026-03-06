package command

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/harness/harness-cli/cmd/cmdutils"
	"github.com/harness/harness-cli/config"

	"github.com/spf13/cobra"
)

// NewGetExecutionCmd creates the get execution command
func NewGetExecutionCmd(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execution EXECUTION_ID",
		Short: "Get execution details",
		Long:  "Retrieves detailed information about a specific pipeline execution",
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

			// Get execution
			response, err := client.GetExecution(context.Background(), orgID, projectID, executionID)
			if err != nil {
				return fmt.Errorf("failed to get execution: %w", err)
			}

			if response.Status != "SUCCESS" {
				return fmt.Errorf("API returned non-success status: %s", response.Status)
			}

			// Print output based on format
			if config.Global.Format == "json" {
				jsonData, err := json.MarshalIndent(response.Data, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal response: %w", err)
				}
				fmt.Println(string(jsonData))
			} else {
				// Print table format
				exec := response.Data.PipelineExecutionSummary
				fmt.Printf("Execution ID: %s\n", exec.PlanExecutionId)
				fmt.Printf("Pipeline: %s (%s)\n", exec.Name, exec.PipelineIdentifier)
				fmt.Printf("Status: %s\n", exec.Status)
				fmt.Printf("Started: %d\n", exec.StartTs)
				if exec.EndTs > 0 {
					fmt.Printf("Ended: %d\n", exec.EndTs)
					duration := exec.EndTs - exec.StartTs
					fmt.Printf("Duration: %d ms\n", duration)
				} else {
					fmt.Printf("Status: Running\n")
				}

				if exec.ExecutionTriggerInfo.TriggerType != "" {
					fmt.Printf("\nTrigger Info:\n")
					fmt.Printf("  Type: %s\n", exec.ExecutionTriggerInfo.TriggerType)
					if exec.ExecutionTriggerInfo.TriggeredBy.Identifier != "" {
						fmt.Printf("  Triggered By: %s\n", exec.ExecutionTriggerInfo.TriggeredBy.Identifier)
					}
				}

				if len(exec.StageNames) > 0 {
					fmt.Printf("\nStages:\n")
					for _, stage := range exec.StageNames {
						fmt.Printf("  - %s\n", stage)
					}
				}
			}

			return nil
		},
	}

	return cmd
}
