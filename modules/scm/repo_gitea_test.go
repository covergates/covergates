// +build gitea

package scm

import (
	"context"
	"net/http"
	"testing"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitea"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

func getGiteaClient() *scm.Client {
	client, _ := gitea.New("http://localhost:3000")
	client.Client = &http.Client{
		Transport: &oauth2.Transport{
			Scheme: oauth2.SchemeBearer,
			Source: &oauth2.Refresher{
				ClientID:     "c8c6a2cc-f948-475c-8663-f420c8fc15ab",
				ClientSecret: "J8YYirhYOZY9a9RepaoORN-8EFcSO-sbwjSGvGo4NwE=",
				Endpoint:     "http://localhost:3000/login/oauth/access_token",
				Source:       oauth2.ContextTokenSource(),
			},
		},
	}
	return client
}

func TestGiteaList(t *testing.T) {

	user := &core.User{
		GiteaToken: "1749a6106454f05f689051c331680c13d78d81b7",
	}
	service := repoService{
		client: getGiteaClient(),
		scm:    core.Gitea,
	}
	repos, _ := service.List(context.Background(), user)
	if len(repos) <= 0 {
		t.Fail()
	}
}

func TestGiteaCreateHook(t *testing.T) {
	service := repoService{
		config: &config.Config{
			Server: config.Server{
				Addr: "http://localhost:8080",
			},
		},
		client: getGiteaClient(),
		scm:    core.Gitea,
	}
	user := &core.User{
		GiteaToken: "1749a6106454f05f689051c331680c13d78d81b7",
	}

	service.CreateHook(context.Background(), user, "gitea/gitea")

}
