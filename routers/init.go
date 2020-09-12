package routers

import (
	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"

	"github.com/covergates/covergates/routers/api"
	"github.com/covergates/covergates/routers/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Routers of server
type Routers struct {
	Config          *config.Config
	Session         core.Session
	LoginMiddleware core.LoginMiddleware
	// service
	SCMService      core.SCMService
	CoverageService core.CoverageService
	ChartService    core.ChartService
	ReportService   core.ReportService
	HookService     core.HookService
	OAuthService    core.OAuthService
	// store
	ReportStore core.ReportStore
	RepoStore   core.RepoStore
}

// RegisterRoutes for Gin engine
func (r *Routers) RegisterRoutes(e *gin.Engine) {
	store := cookie.NewStore([]byte(r.Config.Server.Secret))
	e.Use(sessions.Sessions("codecover", store))
	e.Use(cors.Default())

	webRoute := &web.Router{
		Config:          r.Config,
		LoginMiddleware: r.LoginMiddleware,
		SCMService:      r.SCMService,
		Session:         r.Session,
	}
	apiRoute := &api.Router{
		Config:          r.Config,
		Session:         r.Session,
		CoverageService: r.CoverageService,
		ChartService:    r.ChartService,
		SCMService:      r.SCMService,
		ReportService:   r.ReportService,
		HookService:     r.HookService,
		OAuthService:    r.OAuthService,
		ReportStore:     r.ReportStore,
		RepoStore:       r.RepoStore,
	}
	webRoute.RegisterRoutes(e)
	apiRoute.RegisterRoutes(e)
}
