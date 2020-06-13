package login

import (
	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-login/login"
	"github.com/drone/go-login/login/github"
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
	}
	return middleware
}
