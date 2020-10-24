package login

import (
	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/drone/go-login/login"
	"github.com/drone/go-login/login/gitea"
	"github.com/drone/go-login/login/github"
	"github.com/drone/go-login/login/gitlab"
)

type middleware struct {
	config *config.Config
}

// NewMiddleware of login
func NewMiddleware(config *config.Config) core.LoginMiddleware {
	return &middleware{
		config: config,
	}
}

func (m *middleware) Handler(scm core.SCMProvider) login.Middleware {
	var middleware login.Middleware
	switch scm {
	case core.Github:
		middleware = &github.Config{
			ClientID:     m.config.Github.ClientID,
			ClientSecret: m.config.Github.ClientSecret,
			Server:       m.config.Github.Server,
			Scope:        m.config.Github.Scope,
			Client:       BasicClient(m.config.Github.SkipVerity),
		}
	case core.Gitea:
		middleware = &gitea.Config{
			ClientID:     m.config.Gitea.ClientID,
			ClientSecret: m.config.Gitea.ClientSecret,
			Server:       m.config.Gitea.Server,
			Scope:        m.config.Gitea.Scope,
			RedirectURL:  m.config.Server.URL() + "/login/gitea",
			Client:       BasicClient(m.config.Gitea.SkipVerity),
		}
	case core.GitLab:
		middleware = &gitlab.Config{
			ClientID:     m.config.GitLab.ClientID,
			ClientSecret: m.config.GitLab.ClientSecret,
			RedirectURL:  m.config.Server.URL() + "/login/gitlab",
			Server:       m.config.GitLab.Server,
			Client:       BasicClient(m.config.GitLab.SkipVerity),
			Scope:        m.config.GitLab.Scope,
		}
	case core.Bitbucket:
		middleware = &gitlab.Config{
			ClientID:     m.config.Bitbucket.ClientID,
			ClientSecret: m.config.Bitbucket.ClientSecret,
			RedirectURL:  m.config.Server.URL() + "/login/bitbucket",
			Server:       m.config.GitLab.Server,
			Client:       BasicClient(m.config.GitLab.SkipVerity),
			Scope:        m.config.GitLab.Scope,
		}
	}
	return middleware
}
