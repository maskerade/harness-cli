package pipeline

// PipelineListResponse represents the response from listing pipelines
type PipelineListResponse struct {
	Status string `json:"status"`
	Data   struct {
		TotalPages    int                  `json:"totalPages"`
		TotalItems    int                  `json:"totalItems"`
		PageItemCount int                  `json:"pageItemCount"`
		PageSize      int                  `json:"pageSize"`
		Content       []PipelineSummary    `json:"content"`
		PageIndex     int                  `json:"pageIndex"`
		Empty         bool                 `json:"empty"`
	} `json:"data"`
}

// PipelineSummary represents a pipeline summary in list responses
type PipelineSummary struct {
	Identifier       string   `json:"identifier"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
	CreatedAt        int64    `json:"createdAt"`
	LastUpdatedAt    int64    `json:"lastUpdatedAt"`
	StageNames       []string `json:"stageNames,omitempty"`
	NumOfStages      int      `json:"numOfStages"`
	ExecutionSummaryInfo struct {
		LastExecutionStatus string `json:"lastExecutionStatus,omitempty"`
		LastExecutionTs     int64  `json:"lastExecutionTs,omitempty"`
		NumOfErrors         []int  `json:"numOfErrors,omitempty"`
	} `json:"executionSummaryInfo,omitempty"`
}

// PipelineResponse represents the response from getting a pipeline
type PipelineResponse struct {
	Status string `json:"status"`
	Data   struct {
		Name         string            `json:"name"`
		Identifier   string            `json:"identifier"`
		Description  string            `json:"description,omitempty"`
		Tags         map[string]string `json:"tags,omitempty"`
		YAML         string            `json:"yamlPipeline"`
		CreatedAt    int64             `json:"createdAt"`
		LastUpdatedAt int64            `json:"lastUpdatedAt"`
	} `json:"data"`
}

// ExecutionListResponse represents the response from listing executions
type ExecutionListResponse struct {
	Status string `json:"status"`
	Data   struct {
		TotalPages    int                 `json:"totalPages"`
		TotalItems    int                 `json:"totalItems"`
		PageItemCount int                 `json:"pageItemCount"`
		PageSize      int                 `json:"pageSize"`
		Content       []ExecutionSummary  `json:"content"`
		PageIndex     int                 `json:"pageIndex"`
		Empty         bool                `json:"empty"`
	} `json:"data"`
}

// ExecutionSummary represents an execution summary in list responses
type ExecutionSummary struct {
	PipelineIdentifier string            `json:"pipelineIdentifier"`
	PlanExecutionId    string            `json:"planExecutionId"`
	Name               string            `json:"name"`
	Status             string            `json:"status"`
	StartTs            int64             `json:"startTs"`
	EndTs              int64             `json:"endTs,omitempty"`
	ExecutionTriggerInfo struct {
		TriggerType      string `json:"triggerType"`
		TriggeredBy      struct {
			Uuid       string `json:"uuid"`
			Identifier string `json:"identifier"`
		} `json:"triggeredBy,omitempty"`
	} `json:"executionTriggerInfo,omitempty"`
	ModuleInfo         map[string]interface{} `json:"moduleInfo,omitempty"`
	LayoutNodeMap      map[string]interface{} `json:"layoutNodeMap,omitempty"`
	StageNames         []string               `json:"stageNames,omitempty"`
}

// ExecutionDetailResponse represents the response from getting execution details
type ExecutionDetailResponse struct {
	Status string `json:"status"`
	Data   struct {
		PipelineExecutionSummary ExecutionSummary `json:"pipelineExecutionSummary"`
		ExecutionGraph           struct {
			RootNodeId string                 `json:"rootNodeId"`
			NodeMap    map[string]interface{} `json:"nodeMap"`
		} `json:"executionGraph,omitempty"`
	} `json:"data"`
}

// TriggerResponse represents the response from triggering a pipeline
type TriggerResponse struct {
	Status string `json:"status"`
	Data   struct {
		PlanExecutionId string `json:"planExecutionId"`
		ExecutionYaml   string `json:"executionYaml,omitempty"`
		} `json:"data"`
}

// InputSetListResponse represents the response from listing input sets
type InputSetListResponse struct {
	Status string `json:"status"`
	Data   struct {
		Content []InputSetSummary `json:"content"`
	} `json:"data"`
}

// InputSetSummary represents an input set summary
type InputSetSummary struct {
	Identifier       string            `json:"identifier"`
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	PipelineIdentifier string          `json:"pipelineIdentifier"`
	Tags             map[string]string `json:"tags,omitempty"`
	InputSetType     string            `json:"inputSetType"`
	CreatedAt        int64             `json:"createdAt"`
	LastUpdatedAt    int64             `json:"lastUpdatedAt"`
}

// TriggerListResponse represents the response from listing triggers
type TriggerListResponse struct {
	Status string `json:"status"`
	Data   struct {
		Content []TriggerSummary `json:"content"`
	} `json:"data"`
}

// TriggerSummary represents a trigger summary
type TriggerSummary struct {
	Identifier    string            `json:"identifier"`
	Name          string            `json:"name"`
	Type          string            `json:"type"`
	Enabled       bool              `json:"enabled"`
	Description   string            `json:"description,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	WebhookUrl    string            `json:"webhookUrl,omitempty"`
	CreatedAt     int64             `json:"createdAt"`
	LastUpdatedAt int64             `json:"lastUpdatedAt"`
}
