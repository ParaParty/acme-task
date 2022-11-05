package model

type Config struct {
	Email string `json:"email"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Challenge   Challenge   `json:"challenge"`
	TaskDetails TaskDetails `json:"task-details"`
	Domains     []string    `json:"domains"`
}

type Challenge struct {
	Type       string     `json:"type"`
	Credential Credential `json:"credential"`
}

type TaskDetails struct {
	Type       string     `json:"type"`
	Credential Credential `json:"credential"`

	Services []string `json:"services"`
}

type Credential struct {
	SecretID  string `json:"secret-id"`
	SecretKey string `json:"secret-key"`
}
