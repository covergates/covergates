package main

import (
	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/modules/login"
	"github.com/code-devel-cover/CodeCover/routers"
	"github.com/google/wire"
)

var routerSet = wire.NewSet(
	provideLogin,
	provideRouter,
)

func provideLogin(config *config.Config) core.LoginMiddleware {
	return login.NewMiddleware(config)
}

func provideRouter(
	session core.Session,
	config *config.Config,
	login core.LoginMiddleware,
	// service
	clientService core.SCMClientService,
	repoService core.RepoService,
	userService core.UserService,
	coverageService core.CoverageService,
	// store
	reportStore core.ReportStore,
) *routers.Routers {
	return &routers.Routers{
		LoginMiddleware:  login,
		SCMClientService: clientService,
		Session:          session,
		UserService:      userService,
		RepoService:      repoService,
		Config:           config,
		CoverageService:  coverageService,
		ReportStore:      reportStore,
	}
}
