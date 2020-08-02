package scm

import (
	"context"
	"fmt"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/drone/go-scm/scm"
	"github.com/rs/xid"
)

// RepoService provides repository operations with SCM
type repoService struct {
	config *config.Config
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
		Private:   repo.Private,
	}, nil
}

func (service *repoService) CloneURL(
	ctx context.Context,
	user *core.User,
	name string,
) (string, error) {
	client := service.client
	ctx = withUser(ctx, service.scm, user)
	repo, _, err := client.Repositories.Find(ctx, name)
	if err != nil {
		return "", err
	}
	return repo.Clone, nil
}

func (service *repoService) CreateHook(ctx context.Context, user *core.User, name string) (*core.Hook, error) {
	target := fmt.Sprintf(
		"%s/api/v1/repos/%s/%s/hook",
		service.config.Server.URL(),
		string(service.scm),
		name,
	)
	input := &scm.HookInput{
		Name:       "covergates",
		Target:     target,
		Secret:     service.config.Server.Secret,
		SkipVerify: service.config.Server.SkipVerity,
		Events: scm.HookEvents{
			Push:        true,
			PullRequest: true,
		},
	}
	ctx = withUser(ctx, service.scm, user)
	hook, _, err := service.client.Repositories.CreateHook(ctx, name, input)
	if err != nil {
		return nil, err
	}
	return &core.Hook{
		ID: hook.ID,
	}, nil
}

func (service *repoService) RemoveHook(ctx context.Context, user *core.User, name string, hook *core.Hook) error {
	ctx = withUser(ctx, service.scm, user)
	_, err := service.client.Repositories.DeleteHook(ctx, name, hook.ID)
	return err
}
