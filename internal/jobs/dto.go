package jobs

type CreateJobRequest struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

type JobResponse struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	Retries    int    `json:"retries"`
	MaxRetries int    `json:"max_retries"`
	Error      string `json:"error,omitempty"`
}

type MetricsResponse struct {
	Pending    int `json:"pending"`
	Processing int `json:"processing"`
	Retrying   int `json:"retrying"`
	Completed  int `json:"completed"`
	Failed     int `json:"failed"`
	Total      int `json:"total"`
}
