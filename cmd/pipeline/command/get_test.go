package command

import (
	"io"
	"testing"

	"github.com/harness/harness-cli/util/pipeline"
)

func TestNewGetPipelineCmd(t *testing.T) {
	setupTestConfig()

	tests := []struct {
		name       string
		mockResp   *pipeline.PipelineResponse
		mockErr    error
		args       []string
		wantErr    bool
		errContain string
	}{
		{
			name: "successful get",
			mockResp: &pipeline.PipelineResponse{
				Status: "SUCCESS",
				Data: struct {
					Name          string            `json:"name"`
					Identifier    string            `json:"identifier"`
					Description   string            `json:"description,omitempty"`
					Tags          map[string]string `json:"tags,omitempty"`
					YAML          string            `json:"yamlPipeline"`
					CreatedAt     int64             `json:"createdAt"`
					LastUpdatedAt int64             `json:"lastUpdatedAt"`
				}{
					Name:          "Test Pipeline",
					Identifier:    "test-pipeline",
					Description:   "A test pipeline",
					CreatedAt:     1234567890,
					LastUpdatedAt: 1234567900,
					YAML:          "pipeline:\n  name: Test",
				},
			},
			args:    []string{"test-pipeline"},
			wantErr: false,
		},
		{
			name:       "missing pipeline ID",
			args:       []string{},
			wantErr:    true,
			errContain: "accepts 1 arg(s)",
		},
		{
			name: "with yaml flag",
			mockResp: &pipeline.PipelineResponse{
				Status: "SUCCESS",
				Data: struct {
					Name          string            `json:"name"`
					Identifier    string            `json:"identifier"`
					Description   string            `json:"description,omitempty"`
					Tags          map[string]string `json:"tags,omitempty"`
					YAML          string            `json:"yamlPipeline"`
					CreatedAt     int64             `json:"createdAt"`
					LastUpdatedAt int64             `json:"lastUpdatedAt"`
				}{
					Name:       "Test Pipeline",
					Identifier: "test-pipeline",
					YAML:       "pipeline:\n  name: Test\n  stages: []",
				},
			},
			args:    []string{"test-pipeline", "--yaml"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockPipelineClient{
				getPipelineResp: tt.mockResp,
				err:             tt.mockErr,
			}

			f := newMockFactory(mockClient)
			cmd := NewGetPipelineCmd(f)
			cmd.SetArgs(tt.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			// Verify command structure
			if cmd.Use != "get PIPELINE_ID" {
				t.Errorf("Command.Use = %v, want %v", cmd.Use, "get PIPELINE_ID")
			}

			// Test argument validation only for positional args
			positionalArgs := extractPositionalArgs(tt.args)
			err := cmd.Args(cmd, positionalArgs)
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGetPipelineCmd_Flags(t *testing.T) {
	setupTestConfig()
	mockClient := &mockPipelineClient{}
	f := newMockFactory(mockClient)
	cmd := NewGetPipelineCmd(f)

	// Verify yaml flag exists
	if cmd.Flags().Lookup("yaml") == nil {
		t.Error("yaml flag not found")
	}

	// Verify default value
	yamlFlag := cmd.Flags().Lookup("yaml")
	if yamlFlag.DefValue != "false" {
		t.Errorf("yaml default = %v, want false", yamlFlag.DefValue)
	}
}

// Helper function to extract non-flag arguments
func extractPositionalArgs(args []string) []string {
	var result []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--yaml" {
			continue
		}
		if len(args[i]) > 0 && args[i][0] != '-' {
			result = append(result, args[i])
		}
	}
	return result
}
