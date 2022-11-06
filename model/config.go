package model

type Config struct {
	Tasks []Task     `json:"tasks"`
	Acme  AcmeConfig `json:"acme"`
}

type AcmeConfig struct {
	Email          string            `json:"email"`
	Type           string            `json:"type"`
	HmacEncoded    string            `json:"hmac_encoded"`
	KeyId          string            `json:"key_id"`
	ValidityPeriod string            `json:"validity_period"`
	Details        AcmeConfigDetails `json:"details"`
}

type AcmeConfigDetails struct {
	Project    string            `json:"project"`
	Credential map[string]string `json:"credential"`
}

type Task struct {
	Challenge   Challenge   `json:"challenge"`
	TaskDetails TaskDetails `json:"task_details"`
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
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}
