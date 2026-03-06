package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/harness/harness-cli/config"
)

// Client handles communication with the Harness Pipeline API
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
	accountID  string
}

// NewClient creates a new pipeline API client
func NewClient() *Client {
	return &Client{
		baseURL:    config.Global.APIBaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      config.Global.AuthToken,
		accountID:  config.Global.AccountID,
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// ListPipelines lists all pipelines in an organization/project
func (c *Client) ListPipelines(ctx context.Context, orgID, projectID string, size, page int) (*PipelineListResponse, error) {
	params := url.Values{}
	params.Add("accountIdentifier", c.accountID)
	params.Add("orgIdentifier", orgID)
	params.Add("projectIdentifier", projectID)
	if size > 0 {
		params.Add("size", fmt.Sprintf("%d", size))
	}
	if page > 0 {
		params.Add("page", fmt.Sprintf("%d", page))
	}

	path := fmt.Sprintf("/pipeline/api/pipelines/list?%s", params.Encode())
	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PipelineListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetPipeline gets details of a specific pipeline
func (c *Client) GetPipeline(ctx context.Context, orgID, projectID, pipelineID string) (*PipelineResponse, error) {
	params := url.Values{}
	params.Add("accountIdentifier", c.accountID)
	params.Add("orgIdentifier", orgID)
	params.Add("projectIdentifier", projectID)

	path := fmt.Sprintf("/pipeline/api/pipelines/%s?%s", pipelineID, params.Encode())
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PipelineResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ListExecutions lists pipeline executions
func (c *Client) ListExecutions(ctx context.Context, orgID, projectID, pipelineID, status string, size, page int) (*ExecutionListResponse, error) {
	params := url.Values{}
	params.Add("accountIdentifier", c.accountID)
	params.Add("orgIdentifier", orgID)
	params.Add("projectIdentifier", projectID)
	if pipelineID != "" {
		params.Add("pipelineIdentifier", pipelineID)
	}
	if status != "" {
		params.Add("status", status)
	}
	if size > 0 {
		params.Add("size", fmt.Sprintf("%d", size))
	}
	if page > 0 {
		params.Add("page", fmt.Sprintf("%d", page))
	}

	path := fmt.Sprintf("/pipeline/api/pipelines/execution/summary?%s", params.Encode())
	resp, err := c.doRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ExecutionListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetExecution gets details of a specific execution
func (c *Client) GetExecution(ctx context.Context, orgID, projectID, planExecutionID string) (*ExecutionDetailResponse, error) {
	params := url.Values{}
	params.Add("accountIdentifier", c.accountID)
	params.Add("orgIdentifier", orgID)
	params.Add("projectIdentifier", projectID)

	path := fmt.Sprintf("/pipeline/api/pipelines/execution/%s?%s", planExecutionID, params.Encode())
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ExecutionDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetExecutionLogs downloads logs for a pipeline execution
func (c *Client) GetExecutionLogs(ctx context.Context, orgID, projectID, planExecutionID string) ([]byte, error) {
	params := url.Values{}
	params.Add("accountIdentifier", c.accountID)
	params.Add("orgIdentifier", orgID)
	params.Add("projectIdentifier", projectID)

	path := fmt.Sprintf("/pipeline/api/pipelines/execution/%s/logs?%s", planExecutionID, params.Encode())
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs: %w", err)
	}

	return logs, nil
}

// TriggerPipeline triggers a pipeline execution
func (c *Client) TriggerPipeline(ctx context.Context, orgID, projectID, pipelineID string, inputSetRefs []string, runtimeInputYAML string) (*TriggerResponse, error) {
	params := url.Values{}
	params.Add("accountIdentifier", c.accountID)
	params.Add("orgIdentifier", orgID)
	params.Add("projectIdentifier", projectID)

	reqBody := map[string]interface{}{
		"inputSetReferences": inputSetRefs,
	}
	if runtimeInputYAML != "" {
		reqBody["runtimeInputYaml"] = runtimeInputYAML
	}

	path := fmt.Sprintf("/pipeline/api/pipelines/execution/%s?%s", pipelineID, params.Encode())
	resp, err := c.doRequest(ctx, "POST", path, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TriggerResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ListInputSets lists input sets for a pipeline
func (c *Client) ListInputSets(ctx context.Context, orgID, projectID, pipelineID string) (*InputSetListResponse, error) {
	params := url.Values{}
	params.Add("accountIdentifier", c.accountID)
	params.Add("orgIdentifier", orgID)
	params.Add("projectIdentifier", projectID)
	params.Add("pipelineIdentifier", pipelineID)

	path := fmt.Sprintf("/pipeline/api/inputSets?%s", params.Encode())
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result InputSetListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ListTriggers lists triggers for a pipeline
func (c *Client) ListTriggers(ctx context.Context, orgID, projectID, pipelineID string) (*TriggerListResponse, error) {
	params := url.Values{}
	params.Add("accountIdentifier", c.accountID)
	params.Add("orgIdentifier", orgID)
	params.Add("projectIdentifier", projectID)
	params.Add("targetIdentifier", pipelineID)

	path := fmt.Sprintf("/pipeline/api/triggers?%s", params.Encode())
	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TriggerListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
