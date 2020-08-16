package scm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/covergates/covergates/core"
	"github.com/drone/go-scm/scm"
)

type gitService struct {
	git       core.Git
	scm       core.SCMProvider
	scmClient *scm.Client
}

type giteaCommit struct {
	Sha       string           `json:"sha"`
	Commit    *giteaRepoCommit `json:"commit"`
	Committer *giteaUser       `json:"committer"`
}

type giteaRepoCommit struct {
	Message   string           `json:"message"`
	Committer *giteaCommitUser `json:"committer"`
}

type giteaCommitUser struct {
	Name string `json:"name"`
}

type giteaUser struct {
	UserName string `json:"username"`
	Avatar   string `json:"avatar_url"`
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

func (service *gitService) ListCommits(ctx context.Context, user *core.User, repo string) ([]*core.Commit, error) {
	if service.scm == core.Gitea {
		return service.listGiteaCommits(ctx, user, repo)
	}
	client := service.scmClient
	ctx = withUser(ctx, service.scm, user)
	commits, _, err := client.Git.ListCommits(ctx, repo, scm.CommitListOptions{Size: 20})
	if err != nil {
		return nil, err
	}
	results := make([]*core.Commit, len(commits))
	for i, commit := range commits {
		results[i] = &core.Commit{
			Sha:             commit.Sha,
			Message:         commit.Message,
			Committer:       commit.Committer.Name,
			CommitterAvater: commit.Committer.Avatar,
		}
	}

	return results, nil
}

func (service *gitService) listGiteaCommits(ctx context.Context, user *core.User, repo string) ([]*core.Commit, error) {
	client := service.scmClient
	ctx = withUser(ctx, service.scm, user)
	res, err := client.Do(ctx, &scm.Request{
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Method: "GET",
		Path:   fmt.Sprintf("api/v1/repos/%s/commits", repo),
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.Status > 300 {
		return nil, errors.New(http.StatusText(res.Status))
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var commits []giteaCommit
	if err := json.Unmarshal(data, &commits); err != nil {
		return nil, err
	}
	results := make([]*core.Commit, len(commits))
	for i, commit := range commits {
		results[i] = &core.Commit{
			Sha:     commit.Sha,
			Message: commit.Commit.Message,
		}
		if commit.Committer != nil {
			results[i].Committer = commit.Committer.UserName
			results[i].CommitterAvater = commit.Committer.Avatar
		} else {
			results[i].Committer = commit.Commit.Committer.Name
		}
	}
	return results, nil
}

// GitRepository clone
func (service *gitService) GitRepository(ctx context.Context, user *core.User, repo string) (core.GitRepository, error) {
	client := service.scmClient
	rs := &repoService{scm: service.scm, client: client}
	token := userToken(service.scm, user)
	ctx = withUser(ctx, service.scm, user)
	url, err := rs.CloneURL(ctx, user, repo)
	if err != nil {
		return nil, err
	}
	return service.git.Clone(ctx, url, token.Token)
}
