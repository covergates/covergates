package repo

import (
	"context"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
	"github.com/rs/xid"
)

// Service provides repository operations with SCM
type Service struct {
	ClientService core.SCMClientService
}

// NewReportID for upload report
func (service *Service) NewReportID(repo *core.Repo) string {
	guid := xid.New()
	return guid.String()
}

// TODO: Need to add test case

// List repositories from SCM
func (service *Service) List(
	ctx context.Context,
	s core.SCMProvider,
	user *core.User,
) ([]*core.Repo, error) {
	client, err := service.ClientService.Client(s)
	ctx = service.ClientService.WithUser(ctx, s, user)
	if err != nil {
		return nil, err
	}
	result, _, err := client.Repositories.List(ctx, scm.ListOptions{Size: 100})
	if err != nil {
		return nil, err
	}
	repositories := make([]*core.Repo, len(result))
	for i, r := range result {
		repositories[i] = &core.Repo{
			Name: r.Namespace,
			URL:  r.Link,
			SCM:  s,
		}
	}
	return repositories, nil
}
