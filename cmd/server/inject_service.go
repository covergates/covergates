package main

import (
	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/modules/charts"
	"github.com/code-devel-cover/CodeCover/modules/git"
	"github.com/code-devel-cover/CodeCover/modules/report"
	"github.com/code-devel-cover/CodeCover/modules/scm"
	"github.com/code-devel-cover/CodeCover/modules/session"
	"github.com/code-devel-cover/CodeCover/service/coverage"
	"github.com/google/wire"
)

var serviceSet = wire.NewSet(
	provideSCMService,
	provideSession,
	provideCoverageService,
	provideChartService,
	provideGit,
	provideReportService,
)

func provideSCMService(
	config *config.Config,
	userStore core.UserStore,
	git core.Git,
) core.SCMService {
	return &scm.Service{
		Config:    config,
		UserStore: userStore,
		Git:       git,
	}
}

func provideSession() core.Session {
	return &session.Session{}
}

func provideCoverageService() core.CoverageService {
	return &coverage.Service{}
}

func provideChartService() core.ChartService {
	return &charts.ChartService{}
}

func provideGit() core.Git {
	return &git.Service{}
}

func provideReportService() core.ReportService {
	return &report.Service{}
}
