package user

import (
	"context"

	"github.com/code-devel-cover/CodeCover/core"
)

type UserService struct {
	ClientService core.SCMClientService
	UserStore     core.UserStore
}

func (u *UserService) Find(ctx context.Context, s core.SCMProvider) (*core.User, error) {
	client := u.ClientService.Client(s)
	user, _, err := client.Users.Find(ctx)
	if err != nil {
		return nil, err
	}
	return u.UserStore.Find(s, user)
}
