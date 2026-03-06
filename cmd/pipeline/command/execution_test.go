package command

import (
	"io"
	"testing"

	"github.com/harness/harness-cli/util/pipeline"
)

func TestNewGetExecutionCmd(t *testing.T) {
	setupTestConfig()

	tests := []struct {
		name       string
		mockResp   *pipeline.ExecutionDetailResponse
		mockErr    error
		args       []string
		wantErr    bool
		errContain string
	}{
		{
			name: "successful get execution",
			mockResp: &pipeline.ExecutionDetailResponse{
				Status: "SUCCESS",
				Data: struct {
					PipelineExecutionSummary pipeline.ExecutionSummary `json:"pipelineExecutionSummary"`
					ExecutionGraph           struct {
						RootNodeId string                 `json:"rootNodeId"`
						NodeMap    map[string]interface{} `json:"nodeMap"`
					} `json:"executionGraph,omitempty"`
				}{
					PipelineExecutionSummary: pipeline.ExecutionSummary{
						PipelineIdentifier: "test-pipeline",
						PlanExecutionId:    "exec-123",
						Name:               "Test Execution",
						Status:             "SUCCESS",
						StartTs:            1234567890,
						EndTs:              1234567900,
					},
				},
			},
			args:    []string{"exec-123"},
			wantErr: false,
		},
		{
			name:       "missing execution ID",
			args:       []string{},
			wantErr:    true,
			errContain: "accepts 1 arg(s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockPipelineClient{
				getExecutionResp: tt.mockResp,
				err:              tt.mockErr,
			}

			f := newMockFactory(mockClient)
			cmd := NewGetExecutionCmd(f)
			cmd.SetArgs(tt.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			// Verify command structure
			if cmd.Use != "execution EXECUTION_ID" {
				t.Errorf("Command.Use = %v, want %v", cmd.Use, "execution EXECUTION_ID")
			}

			// Test argument validation
			err := cmd.Args(cmd, tt.args)
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestNewGetLogsCmd(t *testing.T) {
	setupTestConfig()

	tests := []struct {
		name       string
		mockLogs   []byte
		mockErr    error
		args       []string
		wantErr    bool
		errContain string
	}{
		{
			name:     "successful log download",
			mockLogs: []byte("test log content\nline 2\nline 3"),
			args:     []string{"exec-123"},
			wantErr:  false,
		},
		{
			name:       "missing execution ID",
			args:       []string{},
			wantErr:    true,
			errContain: "accepts 1 arg(s)",
		},
		{
			name:     "with custom output file",
			mockLogs: []byte("test logs"),
			args:     []string{"exec-123", "--output", "custom-logs.txt"},
			wantErr:  false,
		},
		{
			name:     "with output directory",
			mockLogs: []byte("test logs"),
			args:     []string{"exec-123", "--output-dir", "./logs"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockPipelineClient{
				getLogsResp: tt.mockLogs,
				err:         tt.mockErr,
			}

			f := newMockFactory(mockClient)
			cmd := NewGetLogsCmd(f)
			cmd.SetArgs(tt.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			// Verify command structure
			if cmd.Use != "logs EXECUTION_ID" {
				t.Errorf("Command.Use = %v, want %v", cmd.Use, "logs EXECUTION_ID")
			}

			// Test argument validation
			execIDArgs := extractExecutionIDArgs(tt.args)
			err := cmd.Args(cmd, execIDArgs)
			if tt.wantErr && tt.errContain == "accepts 1 arg(s)" && err == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}

func TestGetLogsCmd_Flags(t *testing.T) {
	setupTestConfig()
	mockClient := &mockPipelineClient{}
	f := newMockFactory(mockClient)
	cmd := NewGetLogsCmd(f)

	// Verify flags exist
	if cmd.Flags().Lookup("output-dir") == nil {
		t.Error("output-dir flag not found")
	}
	if cmd.Flags().Lookup("output") == nil {
		t.Error("output flag not found")
	}

	// Verify short flag alias for output
	outputFlag := cmd.Flags().ShorthandLookup("o")
	if outputFlag == nil {
		t.Error("output short flag -o not found")
	}
}

// Helper function to extract non-flag arguments for execution ID
func extractExecutionIDArgs(args []string) []string {
	var result []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--output" || args[i] == "-o" || args[i] == "--output-dir" {
			i++ // Skip the value
			continue
		}
		if len(args[i]) > 0 && args[i][0] != '-' {
			result = append(result, args[i])
		}
	}
	return result
}
