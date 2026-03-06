package command

import (
	"context"

	"github.com/harness/harness-cli/cmd/cmdutils"
	"github.com/harness/harness-cli/config"
	"github.com/harness/harness-cli/util/pipeline"
)

// mockPipelineClient implements a mock pipeline client for testing
type mockPipelineClient struct {
	listPipelinesResp   *pipeline.PipelineListResponse
	getPipelineResp     *pipeline.PipelineResponse
	listExecutionsResp  *pipeline.ExecutionListResponse
	getExecutionResp    *pipeline.ExecutionDetailResponse
	getLogsResp         []byte
	triggerResp         *pipeline.TriggerResponse
	listInputSetsResp   *pipeline.InputSetListResponse
	listTriggersResp    *pipeline.TriggerListResponse
	err                 error
}

func (m *mockPipelineClient) ListPipelines(ctx context.Context, orgID, projectID string, size, page int) (*pipeline.PipelineListResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.listPipelinesResp, nil
}

func (m *mockPipelineClient) GetPipeline(ctx context.Context, orgID, projectID, pipelineID string) (*pipeline.PipelineResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.getPipelineResp, nil
}

func (m *mockPipelineClient) ListExecutions(ctx context.Context, orgID, projectID, pipelineID, status string, size, page int) (*pipeline.ExecutionListResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.listExecutionsResp, nil
}

func (m *mockPipelineClient) GetExecution(ctx context.Context, orgID, projectID, planExecutionID string) (*pipeline.ExecutionDetailResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.getExecutionResp, nil
}

func (m *mockPipelineClient) GetExecutionLogs(ctx context.Context, orgID, projectID, planExecutionID string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.getLogsResp, nil
}

func (m *mockPipelineClient) TriggerPipeline(ctx context.Context, orgID, projectID, pipelineID string, inputSetRefs []string, runtimeInputYAML string) (*pipeline.TriggerResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.triggerResp, nil
}

func (m *mockPipelineClient) ListInputSets(ctx context.Context, orgID, projectID, pipelineID string) (*pipeline.InputSetListResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.listInputSetsResp, nil
}

func (m *mockPipelineClient) ListTriggers(ctx context.Context, orgID, projectID, pipelineID string) (*pipeline.TriggerListResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.listTriggersResp, nil
}

func setupTestConfig() {
	config.Global = config.GlobalFlags{
		AccountID: "test-account",
		OrgID:     "test-org",
		ProjectID: "test-project",
		Format:    "table",
	}
}

func newMockFactory(mockClient *mockPipelineClient) *cmdutils.Factory {
	return &cmdutils.Factory{
		PipelineClient: func() *pipeline.Client {
			// Return the mock client by casting it to the interface
			// This is a workaround since we can't return the interface directly
			return (*pipeline.Client)(nil) // We'll handle this in the command tests
		},
	}
}
