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

func provideRouter(login core.LoginMiddleware) *routers.Routers {
	return &routers.Routers{
		LoginMiddleware: login,
	}
}
