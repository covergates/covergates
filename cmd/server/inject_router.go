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
	scmService core.SCMService,
	coverageService core.CoverageService,
	chartService core.ChartService,
	// store
	reportStore core.ReportStore,
	repoStore core.RepoStore,
) *routers.Routers {
	return &routers.Routers{
		Config:          config,
		Session:         session,
		LoginMiddleware: login,
		SCMService:      scmService,
		CoverageService: coverageService,
		ChartService:    chartService,
		ReportStore:     reportStore,
		RepoStore:       repoStore,
	}
}
