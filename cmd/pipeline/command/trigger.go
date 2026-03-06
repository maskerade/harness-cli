package command

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/harness/harness-cli/cmd/cmdutils"
	"github.com/harness/harness-cli/config"

	"github.com/spf13/cobra"
)

// NewTriggerPipelineCmd creates the trigger command
func NewTriggerPipelineCmd(f *cmdutils.Factory) *cobra.Command {
	var inputSetRefs []string
	var runtimeInputFile string
	var confirm bool

	cmd := &cobra.Command{
		Use:   "trigger PIPELINE_ID",
		Short: "Trigger a pipeline execution",
		Long:  "Triggers a pipeline execution with optional input sets and runtime inputs",
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

			// Require confirmation for write operations
			if !confirm {
				fmt.Println("⚠️  Pipeline trigger requires confirmation.")
				fmt.Printf("Pipeline: %s\n", pipelineID)
				fmt.Printf("Organization: %s\n", orgID)
				fmt.Printf("Project: %s\n", projectID)
				if len(inputSetRefs) > 0 {
					fmt.Printf("Input Sets: %s\n", strings.Join(inputSetRefs, ", "))
				}
				fmt.Println("\nAdd --confirm flag to proceed with the trigger")
				return nil
			}

			// Read runtime input file if provided
			var runtimeInputYAML string
			if runtimeInputFile != "" {
				data, err := os.ReadFile(runtimeInputFile)
				if err != nil {
					return fmt.Errorf("failed to read runtime input file: %w", err)
				}
				runtimeInputYAML = string(data)
			}

			// Create pipeline client
			client := f.PipelineClient()

			// Trigger pipeline
			response, err := client.TriggerPipeline(context.Background(), orgID, projectID, pipelineID, inputSetRefs, runtimeInputYAML)
			if err != nil {
				return fmt.Errorf("failed to trigger pipeline: %w", err)
			}

			if response.Status != "SUCCESS" {
				return fmt.Errorf("API returned non-success status: %s", response.Status)
			}

			fmt.Printf("✓ Pipeline triggered successfully\n")
			fmt.Printf("Execution ID: %s\n", response.Data.PlanExecutionId)
			fmt.Printf("\nUse 'hc pipeline execution %s' to check status\n", response.Data.PlanExecutionId)

			return nil
		},
	}

	cmd.Flags().StringSliceVar(&inputSetRefs, "input-set", []string{}, "input set references (can be specified multiple times)")
	cmd.Flags().StringVar(&runtimeInputFile, "runtime-input", "", "path to runtime input YAML file")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "confirm pipeline trigger (required for execution)")

	return cmd
}
