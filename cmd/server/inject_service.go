package main

import (
	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/modules/scm"
	"github.com/code-devel-cover/CodeCover/modules/session"
	"github.com/code-devel-cover/CodeCover/modules/user"
	"github.com/google/wire"
)

var serviceSet = wire.NewSet(
	provideSCMClientService,
	provideUserService,
	provideSession,
)

func provideSCMClientService(config *config.Config) core.SCMClientService {
	return scm.NewSCMClientService(config)
}

func provideUserService(userStore core.UserStore, client core.SCMClientService) core.UserService {
	return &user.Service{
		UserStore:     userStore,
		ClientService: client,
	}
}

func provideSession() core.Session {
	return &session.Session{}
}
