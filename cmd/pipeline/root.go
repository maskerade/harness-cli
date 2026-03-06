package pipeline

import (
	"github.com/harness/harness-cli/cmd/cmdutils"
	"github.com/harness/harness-cli/cmd/pipeline/command"

	"github.com/spf13/cobra"
)

func GetRootCmd(f *cmdutils.Factory) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "pipeline",
		Aliases: []string{"pipe"},
		Short:   "Manage Harness Pipelines",
		Long:    `Commands to manage Harness CI/CD Pipelines and their executions`,
	}

	// Add subcommands
	rootCmd.AddCommand(command.NewListPipelineCmd(f))
	rootCmd.AddCommand(command.NewGetPipelineCmd(f))
	rootCmd.AddCommand(command.NewListExecutionsCmd(f))
	rootCmd.AddCommand(command.NewGetExecutionCmd(f))
	rootCmd.AddCommand(command.NewGetLogsCmd(f))
	rootCmd.AddCommand(command.NewTriggerPipelineCmd(f))
	rootCmd.AddCommand(command.NewListInputSetsCmd(f))
	rootCmd.AddCommand(command.NewListTriggersCmd(f))

	return rootCmd}
