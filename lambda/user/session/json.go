package session

// LambdaSession is a session for a logged in Lambda user.
type LambdaSession struct {
	// sessionID is the ID associated with the user session. If empty, it means that a new
	// session should be created and attached to the request when SaveToContext is called.
	sessionID string `json:"-"`
	// APIKey is the API key of the user that the session belongs to.
	APIKey string `json:"api_key"`
}

// NewSession creates a new Lambda session that will get assigned a session ID when persisted.
func NewSession(apiKey string) LambdaSession {
	return LambdaSession{"", apiKey}
}
