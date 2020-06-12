package web

import (
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/gin-gonic/gin"
)
type WebRouter struct {
	LoginMiddleware  core.LoginMiddleware
	SCMClientService core.SCMClientService
	UserService      core.UserService
	Session          core.Session
}

func (r *WebRouter) RegisterRoutes(e *gin.Engine) {
	{
		g := e.Group("/login")
		g.Any("/github",
			MiddlewareLogin(core.Github, r.LoginMiddleware),
			HandleLogin(
				core.Github,
				r.SCMClientService,
				r.UserService,
				r.Session,
			),
		)
	}
	e.Any("/logout", handleLogout(r.Session))
	e.NoRoute(HandleIndex)
}
