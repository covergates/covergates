package core

import (
	"context"

	"github.com/drone/go-scm/scm"
)

//go:generate mockgen -package mock -destination ../mock/scm_mock.go . SCMClientService

type SCMClientService interface {
	Client(scm SCMProvider) (*scm.Client, error)
	WithUser(ctx context.Context, scm SCMProvider, usr *User) context.Context
}
