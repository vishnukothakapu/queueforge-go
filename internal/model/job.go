package model

type Job struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	Data       string `json:"data"`
	Retries    int    `json:"retries"`
	MaxRetries int    `json:"max_retries"`
}
