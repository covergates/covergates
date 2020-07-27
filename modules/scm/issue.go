package scm

import (
	"context"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
)

type issueService struct {
	client *scm.Client
	scm    core.SCMProvider
}

func (service *issueService) CreateComment(
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

func (service *issueService) RemoveComment(
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
