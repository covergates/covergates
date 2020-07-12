package scm

import (
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
)

type client struct {
	scm       core.SCMProvider
	scmClient *scm.Client
	git       core.Git
	userStore core.UserStore
}

func (c *client) Repositories() core.RepoService {
	return &repoService{
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
