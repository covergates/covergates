package scm

import (
	"context"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
	"github.com/rs/xid"
)

// RepoService provides repository operations with SCM
type repoService struct {
	client *scm.Client
	scm    core.SCMProvider
}

// NewReportID for upload report
func (service *repoService) NewReportID(repo *core.Repo) string {
	guid := xid.New()
	return guid.String()
}

// List repositories from SCM
func (service *repoService) List(
	ctx context.Context,
	user *core.User,
) ([]*core.Repo, error) {
	client := service.client
	ctx = withUser(ctx, service.scm, user)
	result, _, err := client.Repositories.List(ctx, scm.ListOptions{Size: 100})
	if err != nil {
		return nil, err
	}
	repositories := make([]*core.Repo, len(result))
	for i, r := range result {
		repositories[i] = &core.Repo{
			NameSpace: r.Namespace,
			Name:      r.Name,
			URL:       r.Link,
			SCM:       service.scm,
			Branch:    r.Branch,
		}
	}
	return repositories, nil
}

// Find repository by it's name (namespace/name)
func (service *repoService) Find(
	ctx context.Context,
	user *core.User,
	name string,
) (*core.Repo, error) {
	client := service.client
	ctx = withUser(ctx, service.scm, user)
	repo, _, err := client.Repositories.Find(ctx, name)
	if err != nil {
		return nil, err
	}
	return &core.Repo{
		Name:      repo.Name,
		NameSpace: repo.Namespace,
		SCM:       core.SCMProvider(service.scm),
		URL:       repo.Link,
		Branch:    repo.Branch,
	}, nil
}
