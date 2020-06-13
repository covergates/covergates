package core

import (
	"context"

	"github.com/drone/go-scm/scm"
)

type SCMClientService interface {
	Client(scm SCMProvider) (*scm.Client, error)
	WithUser(ctx context.Context, scm SCMProvider, usr *User) context.Context
}
