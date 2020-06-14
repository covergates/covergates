package repo

import (
	"context"
	"net/http"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/mock"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitea"
	"github.com/drone/go-scm/scm/transport/oauth2"
	"github.com/golang/mock/gomock"
)

func TestNewReportID(t *testing.T) {
	s := &Service{}
	reportID1 := s.NewReportID(&core.Repo{
		URL: "http://repo",
	})
	reportID2 := s.NewReportID(&core.Repo{
		URL: "http://repo",
	})
	if reportID1 == reportID2 {
		t.Logf("%s %s", reportID1, reportID2)
		t.Fail()
	}
}

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	client, _ := gitea.New("http://localhost:3000")
	client.Client = &http.Client{
		Transport: &oauth2.Transport{
			Scheme: oauth2.SchemeBearer,
			Source: &oauth2.Refresher{
				Source: oauth2.ContextTokenSource(),
			},
		},
	}
	ctx := context.Background()
	ctx = scm.WithContext(ctx, &scm.Token{
		Token: "1749a6106454f05f689051c331680c13d78d81b7",
	})
	user := &core.User{
		GiteaToken: "1749a6106454f05f689051c331680c13d78d81b7",
	}
	clientService := mock.NewMockSCMClientService(ctrl)
	clientService.EXPECT().Client(gomock.Eq(core.Gitea)).Return(client, nil)
	clientService.EXPECT().WithUser(
		gomock.Any(),
		gomock.Eq(core.Gitea),
		gomock.Eq(user),
	).Return(ctx)

	service := &Service{
		ClientService: clientService,
	}
	repos, _ := service.List(context.Background(), core.Gitea, user)
	if len(repos) < 1 || repos[0].Name != "repo1" || repos[0].NameSpace != "gitea" {
		t.Fail()
	}
}
