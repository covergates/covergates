package scm

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/bitbucket"
	"github.com/drone/go-scm/scm/driver/gitea"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/driver/gitlab"
	"github.com/drone/go-scm/scm/transport/oauth2"
	log "github.com/sirupsen/logrus"
)

type errClientNotFound struct {
	scm core.SCMProvider
}

func (e *errClientNotFound) Error() string {
	return fmt.Sprintf("%s client not found", e.scm)
}

// Service of SCM
type Service struct {
	Config    *config.Config
	Git       core.Git
	UserStore core.UserStore
}

func transport(insecure bool) http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
		},
	}
}

func userToken(s core.SCMProvider, usr *core.User) *scm.Token {
	var token *scm.Token
	switch s {
	case core.Github:
		token = &scm.Token{
			Token: usr.GithubToken,
		}
	case core.Gitea:
		token = &scm.Token{
			Token:   usr.GiteaToken,
			Refresh: usr.GiteaRefresh,
			Expires: usr.GiteaExpire,
		}
	case core.GitLab:
		token = &scm.Token{
			Token:   usr.GitLabToken,
			Refresh: usr.GitLabRefresh,
			Expires: usr.GitLabExpire,
		}
	case core.Bitbucket:
		token = &scm.Token{
			Token:   usr.BitbucketToken,
			Refresh: usr.BitbucketRefresh,
			Expires: usr.BitbucketExpire,
		}
	default:
		log.Warningf("%s is not supported", s)
	}
	return token
}

func withUser(
	ctx context.Context,
	s core.SCMProvider,
	usr *core.User,
) context.Context {
	return context.WithValue(ctx, scm.TokenKey{}, userToken(s, usr))
}

// Client to access SCM API
func (service *Service) Client(s core.SCMProvider) (core.Client, error) {
	scmClient, err := scmClient(s, service.Config)
	if err != nil {
		return nil, err
	}
	return &client{
		scm:       s,
		config:    service.Config,
		scmClient: scmClient,
		userStore: service.UserStore,
		git:       service.Git,
	}, nil
}

func scmClient(s core.SCMProvider, config *config.Config) (*scm.Client, error) {
	var client *scm.Client
	var err error
	switch s {
	case core.Github:
		client, err = github.New(config.Github.APIServer)
		client.Client = &http.Client{
			Transport: &oauth2.Transport{
				Source: oauth2.ContextTokenSource(),
				Base:   transport(config.Github.SkipVerity),
			},
		}
	case core.Gitea:
		client, err = gitea.New(config.Gitea.Server)
		client.Client = &http.Client{
			Transport: &oauth2.Transport{
				Scheme: oauth2.SchemeBearer,
				Source: &oauth2.Refresher{
					ClientID:     config.Gitea.ClientID,
					ClientSecret: config.Gitea.ClientSecret,
					Endpoint:     strings.TrimPrefix(config.Gitea.Server, "/") + "/login/oauth/access_token",
					Source:       oauth2.ContextTokenSource(),
				},
			},
		}
	case core.GitLab:
		client, err = gitlab.New(config.GitLab.Server)
		client.Client = &http.Client{
			Transport: &oauth2.Transport{
				Source: oauth2.ContextTokenSource(),
				Base: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: config.GitLab.SkipVerity,
					},
				},
			},
		}
	case core.Bitbucket:
		client, err = bitbucket.New(config.Bitbucket.Server)
		client.Client = &http.Client{
			Transport: &oauth2.Transport{
				Source: oauth2.ContextTokenSource(),
				Base: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: config.Bitbucket.SkipVerity,
					},
				},
			},
		}
	default:
		log.Debug("scm not supported")
		return nil, &errClientNotFound{s}
	}
	return client, err
}
