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
// If the user does not exist in storage (database), it returns error even
// the user found in the given SCM.
func (service *userService) Find(ctx context.Context, token *core.Token) (*core.User, error) {
	user, err := service.find(ctx, token)
	if err != nil {
		return nil, err
	}
	return service.store.Find(service.scm, user)
}

func (service *userService) find(ctx context.Context, token *core.Token) (*scm.User, error) {
	client := service.client
	ctx = scm.WithContext(ctx, &scm.Token{
		Token:   token.Token,
		Refresh: token.Refresh,
		Expires: token.Expires,
	})
	user, _, err := client.Users.Find(ctx)
	return user, err
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
	if err := service.store.Create(service.scm, user, token); err != nil {
		return nil, err
	}
	return service.store.Find(service.scm, user)
}

func (service *userService) Bind(ctx context.Context, user *core.User, token *core.Token) (*core.User, error) {
	scmUser, err := service.find(ctx, token)
	if err != nil {
		return user, err
	}
	return service.store.Bind(service.scm, user, scmUser, token)
}
