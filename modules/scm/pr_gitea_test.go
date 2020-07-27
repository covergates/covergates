// +build gitea

package scm

import (
	"context"
	"os"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
)

type mockIssueService struct {
	scm.IssueService
	called bool
}

func (s *mockIssueService) CreateComment(ctx context.Context, repo string, number int, input *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	s.called = true
	return &scm.Comment{ID: 1}, nil, nil
}

func TestGiteaIssueCreate(t *testing.T) {
	mockService := &mockIssueService{}
	service := &prService{
		client: &scm.Client{
			Issues: mockService,
		},
		scm: core.Gitea,
	}

	user := &core.User{
		GiteaToken: os.Getenv("GITEA_SECRET"),
	}

	id, err := service.CreateComment(context.Background(), user, "gitea/JSON", 1, "test")
	if err != nil {
		t.Fatal(err)
	}
	if id <= 0 {
		t.Fail()
	}
	if !mockService.called {
		t.Fail()
	}
}
