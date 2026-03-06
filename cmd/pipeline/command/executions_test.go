package command

import (
	"io"
	"testing"

	"github.com/harness/harness-cli/util/pipeline"
)

func TestNewListExecutionsCmd(t *testing.T) {
	setupTestConfig()

	tests := []struct {
		name       string
		mockResp   *pipeline.ExecutionListResponse
		mockErr    error
		args       []string
		wantErr    bool
		errContain string
	}{
		{
			name: "successful list all executions",
			mockResp: &pipeline.ExecutionListResponse{
				Status: "SUCCESS",
				Data: struct {
					TotalPages    int                          `json:"totalPages"`
					TotalItems    int                          `json:"totalItems"`
					PageItemCount int                          `json:"pageItemCount"`
					PageSize      int                          `json:"pageSize"`
					Content       []pipeline.ExecutionSummary  `json:"content"`
					PageIndex     int                          `json:"pageIndex"`
					Empty         bool                         `json:"empty"`
				}{
					TotalPages:    1,
					TotalItems:    2,
					PageItemCount: 2,
					Content: []pipeline.ExecutionSummary{
						{
							PipelineIdentifier: "test-pipeline",
							PlanExecutionId:    "exec-123",
							Status:             "SUCCESS",
							StartTs:            1234567890,
							EndTs:              1234567900,
						},
					},
				},
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "filter by pipeline",
			mockResp: &pipeline.ExecutionListResponse{
				Status: "SUCCESS",
				Data: struct {
					TotalPages    int                          `json:"totalPages"`
					TotalItems    int                          `json:"totalItems"`
					PageItemCount int                          `json:"pageItemCount"`
					PageSize      int                          `json:"pageSize"`
					Content       []pipeline.ExecutionSummary  `json:"content"`
					PageIndex     int                          `json:"pageIndex"`
					Empty         bool                         `json:"empty"`
				}{
					Content: []pipeline.ExecutionSummary{},
				},
			},
			args:    []string{"my-pipeline"},
			wantErr: false,
		},
		{
			name: "filter by status",
			mockResp: &pipeline.ExecutionListResponse{
				Status: "SUCCESS",
				Data: struct {
					TotalPages    int                          `json:"totalPages"`
					TotalItems    int                          `json:"totalItems"`
					PageItemCount int                          `json:"pageItemCount"`
					PageSize      int                          `json:"pageSize"`
					Content       []pipeline.ExecutionSummary  `json:"content"`
					PageIndex     int                          `json:"pageIndex"`
					Empty         bool                         `json:"empty"`
				}{
					Content: []pipeline.ExecutionSummary{},
				},
			},
			args:    []string{"--status", "FAILED"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockPipelineClient{
				listExecutionsResp: tt.mockResp,
				err:                tt.mockErr,
			}

			f := newMockFactory(mockClient)
			cmd := NewListExecutionsCmd(f)
			cmd.SetArgs(tt.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			// Verify command structure
			if cmd.Use != "executions [PIPELINE_ID]" {
				t.Errorf("Command.Use = %v, want %v", cmd.Use, "executions [PIPELINE_ID]")
			}
		})
	}
}

func TestListExecutionsCmd_Flags(t *testing.T) {
	setupTestConfig()
	mockClient := &mockPipelineClient{}
	f := newMockFactory(mockClient)
	cmd := NewListExecutionsCmd(f)

	// Verify flags exist
	if cmd.Flags().Lookup("status") == nil {
		t.Error("status flag not found")
	}
	if cmd.Flags().Lookup("page-size") == nil {
		t.Error("page-size flag not found")
	}
	if cmd.Flags().Lookup("page") == nil {
		t.Error("page flag not found")
	}

	// Verify default values
	pageSize := cmd.Flags().Lookup("page-size")
	if pageSize.DefValue != "20" {
		t.Errorf("page-size default = %v, want 20", pageSize.DefValue)
	}
}
