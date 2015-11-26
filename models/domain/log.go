package models

type Log struct {
	ID        string `json:"id,omitempty"`
	Token     string `json:"token"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}
