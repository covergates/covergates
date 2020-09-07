package core

import (
	"net/http"
	"time"
)

// OAuthToken holds OAuth2 token information
type OAuthToken struct {
	Code    string
	Access  string
	Refresh string
	Expires time.Time
	Data    []byte
}

// OAuthService provide OAuth2 protocol
type OAuthService interface {
	Authorize(w http.ResponseWriter, r http.Request) error
	Token(w http.ResponseWriter, r http.Request) error
	Validate(r http.Request) (bool, error)
}

// OAuthStore token information
type OAuthStore interface {
	Create(token *OAuthToken) error
	Find(token *OAuthToken) (*OAuthToken, error)
	Delete(token *OAuthToken) error
}
