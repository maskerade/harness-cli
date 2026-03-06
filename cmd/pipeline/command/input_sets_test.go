package command

import (
	"io"
	"testing"

	"github.com/harness/harness-cli/util/pipeline"
)

func TestNewListInputSetsCmd(t *testing.T) {
	setupTestConfig()

	tests := []struct {
		name       string
		mockResp   *pipeline.InputSetListResponse
		mockErr    error
		args       []string
		wantErr    bool
		errContain string
	}{
		{
			name: "successful list input sets",
			mockResp: &pipeline.InputSetListResponse{
				Status: "SUCCESS",
				Data: struct {
					Content []pipeline.InputSetSummary `json:"content"`
				}{
					Content: []pipeline.InputSetSummary{
						{
							Identifier:         "input-set-1",
							Name:               "Production Config",
							PipelineIdentifier: "test-pipeline",
							InputSetType:       "INPUT_SET",
							CreatedAt:          1234567890,
						},
						{
							Identifier:         "input-set-2",
							Name:               "Staging Config",
							PipelineIdentifier: "test-pipeline",
							InputSetType:       "INPUT_SET",
							CreatedAt:          1234567800,
						},
					},
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
			name: "empty input sets",
			mockResp: &pipeline.InputSetListResponse{
				Status: "SUCCESS",
				Data: struct {
					Content []pipeline.InputSetSummary `json:"content"`
				}{
					Content: []pipeline.InputSetSummary{},
				},
			},
			args:    []string{"test-pipeline"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockPipelineClient{
				listInputSetsResp: tt.mockResp,
				err:               tt.mockErr,
			}

			f := newMockFactory(mockClient)
			cmd := NewListInputSetsCmd(f)
			cmd.SetArgs(tt.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			// Verify command structure
			if cmd.Use != "input-sets PIPELINE_ID" {
				t.Errorf("Command.Use = %v, want %v", cmd.Use, "input-sets PIPELINE_ID")
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

func TestNewListTriggersCmd(t *testing.T) {
	setupTestConfig()

	tests := []struct {
		name       string
		mockResp   *pipeline.TriggerListResponse
		mockErr    error
		args       []string
		wantErr    bool
		errContain string
	}{
		{
			name: "successful list triggers",
			mockResp: &pipeline.TriggerListResponse{
				Status: "SUCCESS",
				Data: struct {
					Content []pipeline.TriggerSummary `json:"content"`
				}{
					Content: []pipeline.TriggerSummary{
						{
							Identifier:  "trigger-1",
							Name:        "Webhook Trigger",
							Type:        "WEBHOOK",
							Enabled:     true,
							Description: "Triggered on PR merge",
							CreatedAt:   1234567890,
						},
						{
							Identifier:  "trigger-2",
							Name:        "Scheduled Trigger",
							Type:        "SCHEDULED",
							Enabled:     false,
							Description: "Nightly build",
							CreatedAt:   1234567800,
						},
					},
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
			name: "empty triggers",
			mockResp: &pipeline.TriggerListResponse{
				Status: "SUCCESS",
				Data: struct {
					Content []pipeline.TriggerSummary `json:"content"`
				}{
					Content: []pipeline.TriggerSummary{},
				},
			},
			args:    []string{"test-pipeline"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockPipelineClient{
				listTriggersResp: tt.mockResp,
				err:              tt.mockErr,
			}

			f := newMockFactory(mockClient)
			cmd := NewListTriggersCmd(f)
			cmd.SetArgs(tt.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			// Verify command structure
			if cmd.Use != "triggers PIPELINE_ID" {
				t.Errorf("Command.Use = %v, want %v", cmd.Use, "triggers PIPELINE_ID")
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
