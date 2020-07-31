package scm

import (
	"context"
	"net/http"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitea"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

func TestClientSeccret(t *testing.T) {
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

	ctx := scm.WithContext(context.Background(), &scm.Token{
		Token: "J8YYirhYOZY9a9RepaoORN-8EFcSO-sbwjSGvGo4NwE=",
	})
	_, _, err := client.Issues.CreateComment(ctx, "gitea/JSON", 1, &scm.CommentInput{
		Body: "test",
	})
	if err != nil {
		// TODO: Update gitea image for testing
		t.Log(err)
	}
}
