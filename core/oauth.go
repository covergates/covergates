package core

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

//go:generate mockgen -package mock -destination ../mock/oauth_mock.go . OAuthService,OAuthStore

// OAuthToken holds OAuth2 token information
type OAuthToken struct {
	ID      uint
	Name    string
	Code    string
	Access  string
	Refresh string
	Expires time.Time
	Owner   *User
	Data    []byte
}

// OAuthService provide OAuth2 protocol
type OAuthService interface {
	CreateToken(ctx context.Context, name string) (*OAuthToken, error)
	DeleteToken(ctx context.Context, token *OAuthToken) error
	ListTokens(ctx context.Context) ([]*OAuthToken, error)
	Validate(r *http.Request) (*User, error)
	WithUser(ctx context.Context, user *User) context.Context
}

// OAuthStore token information
type OAuthStore interface {
	Create(token *OAuthToken) error
	Find(token *OAuthToken) (*OAuthToken, error)
	List(user *User) ([]*OAuthToken, error)
	Delete(token *OAuthToken) error
}
