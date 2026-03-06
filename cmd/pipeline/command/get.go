package command

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/harness/harness-cli/cmd/cmdutils"
	"github.com/harness/harness-cli/config"

	"github.com/spf13/cobra"
)

// NewGetPipelineCmd creates the get command for pipelines
func NewGetPipelineCmd(f *cmdutils.Factory) *cobra.Command {
	var showYAML bool

	cmd := &cobra.Command{
		Use:   "get PIPELINE_ID",
		Short: "Get pipeline details",
		Long:  "Retrieves detailed information about a specific pipeline",
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

			// Get pipeline
			response, err := client.GetPipeline(context.Background(), orgID, projectID, pipelineID)
			if err != nil {
				return fmt.Errorf("failed to get pipeline: %w", err)
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
				fmt.Printf("Pipeline: %s\n", response.Data.Name)
				fmt.Printf("Identifier: %s\n", response.Data.Identifier)
				if response.Data.Description != "" {
					fmt.Printf("Description: %s\n", response.Data.Description)
				}
				fmt.Printf("Created At: %d\n", response.Data.CreatedAt)
				fmt.Printf("Last Updated: %d\n", response.Data.LastUpdatedAt)

				if len(response.Data.Tags) > 0 {
					fmt.Printf("\nTags:\n")
					for key, value := range response.Data.Tags {
						fmt.Printf("  %s: %s\n", key, value)
					}
				}

				if showYAML && response.Data.YAML != "" {
					fmt.Printf("\nPipeline YAML:\n%s\n", response.Data.YAML)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&showYAML, "yaml", false, "show pipeline YAML definition")

	return cmd
}
