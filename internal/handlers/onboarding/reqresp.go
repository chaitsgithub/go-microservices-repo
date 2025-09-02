package handlers

// createAppRequest defines the structure for the JSON payload
// received in the request body when creating a new app.
type createAppRequest struct {
	Name string `json:"name"`
}
