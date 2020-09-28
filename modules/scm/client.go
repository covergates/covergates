package scm

import (
	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/drone/go-scm/scm"
)

type client struct {
	config    *config.Config
	scm       core.SCMProvider
	scmClient *scm.Client
	git       core.Git
	userStore core.UserStore
}

func (c *client) Repositories() core.GitRepoService {
	return &repoService{
		config: c.config,
		client: c.scmClient,
		scm:    c.scm,
	}
}

func (c *client) Users() core.UserService {
	return &userService{
		client: c.scmClient,
		store:  c.userStore,
		scm:    c.scm,
	}
}

func (c *client) Git() core.GitService {
	return &gitService{
		git:       c.git,
		scm:       c.scm,
		scmClient: c.scmClient,
	}
}

func (c *client) Contents() core.ContentService {
	return &contentService{
		scm:    c.scm,
		client: c.scmClient,
		git:    c.git,
	}
}

func (c *client) PullRequests() core.PullRequestService {
	return &prService{
		client: c.scmClient,
		scm:    c.scm,
	}
}

func (c *client) Webhooks() core.WebhookService {
	return &webhookService{
		config: c.config,
		client: c.scmClient,
		scm:    c.scm,
	}
}

func (c *client) Token(user *core.User) core.Token {
	token := userToken(c.scm, user)
	return core.Token{
		Token:   token.Token,
		Expires: token.Expires,
		Refresh: token.Refresh,
	}
}
