package web

import (
	"net/http"

	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/web"
	"github.com/gin-gonic/gin"
)

type WebRouter struct {
	Config          *config.Config
	LoginMiddleware core.LoginMiddleware
	SCMService      core.SCMService
	Session         core.Session
}

func (r *WebRouter) RegisterRoutes(e *gin.Engine) {
	{
		g := e.Group("/login")
		g.Any("/github",
			MiddlewareLogin(core.Github, r.LoginMiddleware),
			HandleLogin(
				core.Github,
				r.SCMService,
				r.Session,
			),
		)
		g.Any("/gitea",
			MiddlewareLogin(core.Gitea, r.LoginMiddleware),
			HandleLogin(
				core.Gitea,
				r.SCMService,
				r.Session,
			),
		)
	}
	e.Any("/logout", handleLogout(r.Session))
	h := gin.WrapH(http.FileServer(web.New()))
	e.GET("/favicon.ico", h)
	e.GET("/js/*filepath", h)
	e.GET("/css/*filepath", h)
	e.GET("/img/*filepath", h)
	e.GET("/fonts/*filepath", h)
	e.NoRoute(HandleIndex(r.Config))
}
