package main

import (
	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/modules/scm"
	"github.com/code-devel-cover/CodeCover/modules/session"
	"github.com/code-devel-cover/CodeCover/service/coverage"
	"github.com/google/wire"
)

var serviceSet = wire.NewSet(
	provideSCMService,
	provideSession,
	provideCoverageService,
)

func provideSCMService(config *config.Config, userStore core.UserStore) core.SCMService {
	return &scm.SCMService{
		Config:    config,
		UserStore: userStore,
	}
}

func provideSession() core.Session {
	return &session.Session{}
}

func provideCoverageService() core.CoverageService {
	return &coverage.CoverageService{}
}
