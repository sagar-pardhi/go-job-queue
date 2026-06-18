package jobs

import (
	"time"
)

type Job struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Payload    map[string]any `json:"payload"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	Retries    int            `json:"retries"`
	MaxRetries int            `json:"max_retries"`
	Error      string         `json:"error,omitempty"`
}
