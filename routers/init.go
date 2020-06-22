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
	SCMService      core.SCMService
	CoverageService core.CoverageService
	ChartService    core.ChartService
	// store
	ReportStore core.ReportStore
	RepoStore   core.RepoStore
}

func (r *Routers) RegisterRoutes(e *gin.Engine) {
	store := cookie.NewStore([]byte(r.Config.Server.Secret))
	e.Use(sessions.Sessions("codecover", store))
	e.Use(cors.Default())

	webRoute := &web.WebRouter{
		Config:          r.Config,
		LoginMiddleware: r.LoginMiddleware,
		SCMService:      r.SCMService,
		Session:         r.Session,
	}
	apiRoute := &api.APIRouter{
		Config:          r.Config,
		Session:         r.Session,
		CoverageService: r.CoverageService,
		ChartService:    r.ChartService,
		SCMService:      r.SCMService,
		ReportStore:     r.ReportStore,
		RepoStore:       r.RepoStore,
	}
	webRoute.RegisterRoutes(e)
	apiRoute.RegisterRoutes(e)
}
