package scm

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/transport/oauth2"
	log "github.com/sirupsen/logrus"
)

type scmClientService struct {
	config *config.Config
}

func NewSCMClientService(config *config.Config) core.SCMClientService {
	return &scmClientService{
		config: config,
	}
}

func transport(insecure bool) http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
		},
	}
}

func (service *scmClientService) WithUser(
	ctx context.Context,
	s core.SCMProvider,
	usr *core.User,
) context.Context {
	var token *scm.Token
	switch s {
	case core.Github:
		token = &scm.Token{
			Token: usr.GithubToken,
		}
	}
	return context.WithValue(ctx, scm.TokenKey{}, token)
}

func (service *scmClientService) Client(s core.SCMProvider) *scm.Client {
	var client *scm.Client
	var err error
	switch s {
	case core.Github:
		client, err = github.New(service.config.Github.APIServer)
		client.Client = &http.Client{
			Transport: &oauth2.Transport{
				Source: oauth2.ContextTokenSource(),
				Base:   transport(service.config.Github.SkipVerity),
			},
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	return client
}
