package login

import (
	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-login/login"
	"github.com/drone/go-login/login/github"
)

type loginMiddleware struct {
	config *config.Config
}

func NewLoginMiddleware(config *config.Config) *loginMiddleware {
	return &loginMiddleware{
		config: config,
	}
}

func (m *loginMiddleware) Handler(scm core.SCMProvider) login.Middleware {
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
