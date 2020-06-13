package routers

import (
	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	_ "github.com/code-devel-cover/CodeCover/models"
	"github.com/code-devel-cover/CodeCover/routers/api"
	"github.com/code-devel-cover/CodeCover/routers/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Routers struct {
	Config          *config.Config
	Session         core.Session
	LoginMiddleware core.LoginMiddleware
	// service
	SCMClientService core.SCMClientService
	UserService      core.UserService
	CoverageService  core.CoverageService
	// store
	ReportStore core.ReportStore
}

func (r *Routers) RegisterRoutes(e *gin.Engine) {
	store := cookie.NewStore([]byte(r.Config.Server.Secret))
	e.Use(sessions.Sessions("codecover", store))
	e.Use(cors.Default())

	webRoute := &web.WebRouter{
		LoginMiddleware: r.LoginMiddleware,
		UserService:     r.UserService,
		Session:         r.Session,
	}
	apiRoute := &api.APIRouter{
		CoverageService: r.CoverageService,
		ReportStore:     r.ReportStore,
	}
	webRoute.RegisterRoutes(e)
	apiRoute.RegisterRoutes(e)
}
