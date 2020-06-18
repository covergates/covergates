package scm

import (
	"context"
	"fmt"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
)

type gitService struct {
	scm       core.SCMProvider
	scmClient *scm.Client
}

func (service *gitService) FindCommit(ctx context.Context, user *core.User, repo *core.Repo) string {
	client := service.scmClient
	ctx = withUser(ctx, service.scm, user)
	ref, _, err := client.Git.FindBranch(
		ctx,
		fmt.Sprintf("%s/%s", repo.NameSpace, repo.Name),
		repo.Branch,
	)
	if err != nil {
		return ""
	}
	return ref.Sha
}
