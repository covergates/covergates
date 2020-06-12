package user

import (
	"context"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

// Service for different SCM OAuth token
type Service struct {
	ClientService core.SCMClientService
	UserStore     core.UserStore
}

// Create new user with SCM user
func (u *Service) Create(ctx context.Context, s core.SCMProvider) error {
	client := u.ClientService.Client(s)
	source := oauth2.ContextTokenSource()
	token, err := source.Token(ctx)
	if err != nil {
		return err
	}
	user, _, err := client.Users.Find(ctx)
	if err != nil {
		return err
	}
	return u.UserStore.Create(s, user, token)
}

// Find user with given SCM token from the contex
func (u *Service) Find(ctx context.Context, s core.SCMProvider) (*core.User, error) {
	client := u.ClientService.Client(s)
	user, _, err := client.Users.Find(ctx)
	if err != nil {
		return nil, err
	}
	return u.UserStore.Find(s, user)
}
