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
	return login.NewLoginMiddleware(config)
}

func provideRouter(
	login core.LoginMiddleware,
	client core.SCMClientService,
	user core.UserService,
	session core.Session,
	config *config.Config,
) *routers.Routers {
	return &routers.Routers{
		LoginMiddleware:  login,
		SCMClientService: client,
		Session:          session,
		UserService:      user,
		Config:           config,
	}
}
