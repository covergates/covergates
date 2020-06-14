package web

import (
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/gin-gonic/gin"
)

type WebRouter struct {
	LoginMiddleware core.LoginMiddleware
	UserService     core.UserService
	Session         core.Session
}

func (r *WebRouter) RegisterRoutes(e *gin.Engine) {
	{
		g := e.Group("/login")
		g.Any("/github",
			MiddlewareLogin(core.Github, r.LoginMiddleware),
			HandleLogin(
				core.Github,
				r.UserService,
				r.Session,
			),
		)
		g.Any("/gitea",
			MiddlewareLogin(core.Gitea, r.LoginMiddleware),
			HandleLogin(
				core.Gitea,
				r.UserService,
				r.Session,
			),
		)
	}
	e.Any("/logout", handleLogout(r.Session))
	e.NoRoute(HandleIndex)
}
