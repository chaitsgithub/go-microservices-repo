package models

// App represents the data model for an application.
// The APIKey is stored as a hashed string for security.
type App struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	APIKey string `json:"api_key,omitempty"`
}
