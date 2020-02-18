package user

// Contains all of the structs that will be serialized as json related to login
// This is converted to json_easyjson.go by easyjson

// LoginResult is a result returned when a user attempts to login
type LoginResult struct {
	// Any errors that may have occurred during the login process
	Errors []string `json:"errors"`

	// API key of the logged in user. Only nonempty if login successful.
	APIKey string `json:"api_key"`
	// Whether the login was a success
	Success bool `json:"success"`
}

// RegisterResult is a result returned when a new user attempts to register
type RegisterResult struct {
	// Any errors that may have occurred during the registration process
	Errors []string `json:"errors"`

	// API key of the newly registered user. Only nonempty if registration successful.
	APIKey string `json:"api_key"`
	// Whether the registration was a success
	Success bool `json:"success"`
}

// SessionResult is info about the session that is returned when queried for.
type SessionResult struct {
	// Any errors that may have occurred when checking the session info
	Errors []string `json:"errors"`

	// User ID of the current session
	UserID int `json:"id"`
	// Username of the current session
	Username string `json:"username"`
	// API key of the current user
	APIKey string `json:"api_key"`
}
