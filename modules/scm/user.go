package scm

import (
	"context"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
)

// Service for different SCM OAuth token
type userService struct {
	client *scm.Client
	scm    core.SCMProvider
	store  core.UserStore
}

// Find user with given SCM token from the contex
func (service *userService) Find(ctx context.Context, token *core.Token) (*core.User, error) {
	client := service.client
	ctx = scm.WithContext(ctx, &scm.Token{
		Token:   token.Token,
		Refresh: token.Refresh,
		Expires: token.Expires,
	})
	user, _, err := client.Users.Find(ctx)
	if err != nil {
		return nil, err
	}
	return service.store.Find(service.scm, user)
}

func (service *userService) Create(ctx context.Context, token *core.Token) (*core.User, error) {
	client := service.client
	scmToken := &scm.Token{
		Token:   token.Token,
		Refresh: token.Refresh,
		Expires: token.Expires,
	}
	ctx = scm.WithContext(ctx, scmToken)
	user, _, err := client.Users.Find(ctx)
	if err != nil {
		return nil, err
	}
	if err := service.store.Create(service.scm, user, scmToken); err != nil {
		return nil, err
	}
	return service.store.Find(service.scm, user)
}
