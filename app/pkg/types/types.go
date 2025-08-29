package types

type WebhookPayload struct {
	Path string `json:"path"`
	Body []byte `json:"body"`
}

// Python Agentに渡すペイロード
type TaskPayload struct {
	Source   string                 `json:"source"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
}

// 分類APIからのレスポンス
type ClassificationResponse struct {
	TaskType string `json:"task_type"`
}

// 実行APIからのレスポンス
type ExecutionResponse struct {
	Result string `json:"result"`
}
