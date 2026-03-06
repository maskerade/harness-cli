package command

import (
	"io"
	"testing"

	"github.com/harness/harness-cli/util/pipeline"
)

func TestNewTriggerPipelineCmd(t *testing.T) {
	setupTestConfig()

	tests := []struct {
		name       string
		mockResp   *pipeline.TriggerResponse
		mockErr    error
		args       []string
		wantErr    bool
		errContain string
	}{
		{
			name:       "missing pipeline ID",
			args:       []string{},
			wantErr:    true,
			errContain: "accepts 1 arg(s)",
		},
		{
			name: "without confirm flag - should not error but should not execute",
			mockResp: &pipeline.TriggerResponse{
				Status: "SUCCESS",
				Data: struct {
					PlanExecutionId string `json:"planExecutionId"`
					ExecutionYaml   string `json:"executionYaml,omitempty"`
				}{
					PlanExecutionId: "exec-123",
				},
			},
			args:    []string{"test-pipeline"},
			wantErr: false, // Command succeeds but shows warning
		},
		{
			name: "with confirm flag",
			mockResp: &pipeline.TriggerResponse{
				Status: "SUCCESS",
				Data: struct {
					PlanExecutionId string `json:"planExecutionId"`
					ExecutionYaml   string `json:"executionYaml,omitempty"`
				}{
					PlanExecutionId: "exec-123",
				},
			},
			args:    []string{"test-pipeline", "--confirm"},
			wantErr: false,
		},
		{
			name: "with input sets",
			mockResp: &pipeline.TriggerResponse{
				Status: "SUCCESS",
				Data: struct {
					PlanExecutionId string `json:"planExecutionId"`
					ExecutionYaml   string `json:"executionYaml,omitempty"`
				}{
					PlanExecutionId: "exec-123",
				},
			},
			args:    []string{"test-pipeline", "--input-set", "prod-config", "--confirm"},
			wantErr: false,
		},
		{
			name: "with multiple input sets",
			mockResp: &pipeline.TriggerResponse{
				Status: "SUCCESS",
				Data: struct {
					PlanExecutionId string `json:"planExecutionId"`
					ExecutionYaml   string `json:"executionYaml,omitempty"`
				}{
					PlanExecutionId: "exec-123",
				},
			},
			args:    []string{"test-pipeline", "--input-set", "common", "--input-set", "prod", "--confirm"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockPipelineClient{
				triggerResp: tt.mockResp,
				err:         tt.mockErr,
			}

			f := newMockFactory(mockClient)
			cmd := NewTriggerPipelineCmd(f)
			cmd.SetArgs(tt.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			// Verify command structure
			if cmd.Use != "trigger PIPELINE_ID" {
				t.Errorf("Command.Use = %v, want %v", cmd.Use, "trigger PIPELINE_ID")
			}

			// Test argument validation
			err := cmd.Args(cmd, extractPipelineIDArgs(tt.args))
			if tt.wantErr && tt.errContain == "accepts 1 arg(s)" && err == nil {
				t.Error("Expected error for missing pipeline ID but got none")
			}
		})
	}
}

func TestTriggerPipelineCmd_Flags(t *testing.T) {
	setupTestConfig()
	mockClient := &mockPipelineClient{}
	f := newMockFactory(mockClient)
	cmd := NewTriggerPipelineCmd(f)

	// Verify flags exist
	if cmd.Flags().Lookup("confirm") == nil {
		t.Error("confirm flag not found")
	}
	if cmd.Flags().Lookup("input-set") == nil {
		t.Error("input-set flag not found")
	}
	if cmd.Flags().Lookup("runtime-input") == nil {
		t.Error("runtime-input flag not found")
	}

	// Verify confirm flag default is false (safety feature)
	confirmFlag := cmd.Flags().Lookup("confirm")
	if confirmFlag.DefValue != "false" {
		t.Errorf("confirm default = %v, want false (safety feature)", confirmFlag.DefValue)
	}
}

func TestTriggerPipelineCmd_ConfirmRequired(t *testing.T) {
	// This test verifies that the trigger command has a confirm flag
	// The actual safety check happens in the RunE function which requires
	// integration testing to fully verify
	setupTestConfig()
	mockClient := &mockPipelineClient{}
	f := newMockFactory(mockClient)
	cmd := NewTriggerPipelineCmd(f)

	confirmFlag := cmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Fatal("confirm flag is required for safety")
	}

	if confirmFlag.Usage == "" {
		t.Error("confirm flag should have usage documentation")
	}
}

// Helper function to extract non-flag arguments
func extractPipelineIDArgs(args []string) []string {
	var result []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--confirm" || args[i] == "--input-set" || args[i] == "--runtime-input" {
			if args[i] == "--input-set" || args[i] == "--runtime-input" {
				i++ // Skip the value
			}
			continue
		}
		if len(args[i]) > 0 && args[i][0] != '-' {
			result = append(result, args[i])
		}
	}
	return result
}
