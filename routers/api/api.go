package api

import (
	"net/url"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/routers/api/repo"
	"github.com/covergates/covergates/routers/api/report"
	"github.com/covergates/covergates/routers/api/request"
	"github.com/covergates/covergates/routers/api/user"
	"github.com/covergates/covergates/routers/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:generate swag init -g ./api.go -d ./ -o ../docs

// @title CodeCover API
// @version 1.0
// @description REST API for CodeCover
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// Router for API
type Router struct {
	Config  *config.Config
	Session core.Session
	// service
	CoverageService core.CoverageService
	ChartService    core.ChartService
	SCMService      core.SCMService
	ReportService   core.ReportService
	HookService     core.HookService
	// store
	ReportStore core.ReportStore
	RepoStore   core.RepoStore
}

func host(addr string) string {
	u, err := url.Parse(addr)
	if err != nil {
		return addr
	}
	return u.Host
}

// RegisterRoutes for API
func (r *Router) RegisterRoutes(e *gin.Engine) {
	docs.SwaggerInfo.Host = host(r.Config.Server.Addr)
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	g := e.Group("/api/v1")
	{
		g := g.Group("/user")
		g.GET("", request.CheckLogin(r.Session), user.HandleGet())
		g.GET("/scm", request.CheckLogin(r.Session), user.HandleGetSCM(r.Config))
		g.POST("", user.HandleCreate())
	}
	{
		g := g.Group("/reports")
		g.POST("/:id", report.HandleUpload(
			r.SCMService,
			r.CoverageService,
			r.RepoStore,
			r.ReportStore,
		))
		g.POST("/:id/comment/:number", report.HandleComment(
			r.Config,
			r.SCMService,
			r.RepoStore,
			r.ReportStore,
			r.ReportService,
		))
		g.GET("/:id", report.HandleGet(r.ReportStore, r.RepoStore, r.SCMService))
		g.GET("/:id/treemap/:commit", report.HandleGetTreeMap(
			r.ReportStore,
			r.RepoStore,
			r.ChartService,
		))
	}
	g.GET("/repos/:scm/:namespace/:name", repo.HandleGet(r.RepoStore))
	{
		g := g.Group("/repos")
		g.Use(request.CheckLogin(r.Session))
		g.GET("", repo.HandleListAll(r.Config, r.SCMService, r.RepoStore))
		g.POST("", repo.HandleCreate(r.RepoStore, r.SCMService))
		g.GET("/:scm", repo.HandleListSCM(r.SCMService, r.RepoStore))
		{
			g := g.Group("/:scm/:namespace/:name")
			g.PATCH("", repo.HandleSync(r.SCMService, r.RepoStore))
			g.GET("/setting", repo.HandleGetSetting(r.RepoStore))
			g.POST("/setting", repo.HandleUpdateSetting(r.RepoStore))
			g.PATCH("/report", repo.HandleReportIDRenew(r.RepoStore, r.SCMService))
			g.GET("/files", repo.HandleGetFiles(r.SCMService))
			g.GET("/content/*path", repo.HandleGetFileContent(r.SCMService))
			g.POST("/hook/create", repo.WithRepo(r.RepoStore), repo.HandleHookCreate(r.HookService))
			g.POST("/hook", repo.WithRepo(r.RepoStore), repo.HandleHook(r.SCMService, r.HookService))
		}
	}
}
