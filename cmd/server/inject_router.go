package main

import (
	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/modules/login"
	"github.com/covergates/covergates/routers"
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
	reportService core.ReportService,
	hookService core.HookService,
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
		ReportService:   reportService,
		HookService:     hookService,
		ReportStore:     reportStore,
		RepoStore:       repoStore,
	}
}
