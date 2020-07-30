package scm

import (
	"context"

	"github.com/covergates/covergates/core"
	"github.com/drone/go-scm/scm"
)

type prService struct {
	client *scm.Client
	scm    core.SCMProvider
}

func (service *prService) CreateComment(
	ctx context.Context,
	user *core.User,
	repo string,
	number int,
	body string,
) (int, error) {
	ctx = withUser(ctx, service.scm, user)
	var comment *scm.Comment
	var err error
	input := &scm.CommentInput{Body: body}
	switch service.scm {
	case core.Gitea:
		comment, _, err = service.client.Issues.CreateComment(ctx, repo, number, input)
	default:
		comment, _, err = service.client.PullRequests.CreateComment(ctx, repo, number, input)
	}
	if err != nil {
		return 0, err
	}
	return comment.ID, nil
}

func (service *prService) RemoveComment(
	ctx context.Context,
	user *core.User,
	repo string,
	number int,
	id int,
) error {
	ctx = withUser(ctx, service.scm, user)
	var err error
	switch service.scm {
	case core.Gitea:
		_, err = service.client.Issues.DeleteComment(ctx, repo, number, id)
	default:
		_, err = service.client.PullRequests.DeleteComment(ctx, repo, number, id)
	}
	return err
}

func (service *prService) Find(ctx context.Context, user *core.User, repo string, number int) (*core.PullRequest, error) {
	ctx = withUser(ctx, service.scm, user)
	pr, _, err := service.client.PullRequests.Find(ctx, repo, number)
	if err != nil {
		return nil, err
	}

	return &core.PullRequest{
		Number: pr.Number,
		Commit: pr.Sha,
		Source: pr.Source,
		Target: pr.Target,
	}, nil
}

func (service *prService) ListChanges(ctx context.Context, user *core.User, repo string, number int) ([]*core.FileChange, error) {
	ctx = withUser(ctx, service.scm, user)
	changes, _, err := service.client.PullRequests.ListChanges(ctx, repo, number, scm.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]*core.FileChange, len(changes))
	for i, change := range changes {
		result[i] = &core.FileChange{
			Path:    change.Path,
			Added:   change.Added,
			Deleted: change.Deleted,
			Renamed: change.Renamed,
		}
	}
	return result, nil
}
