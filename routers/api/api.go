package api

import (
	"net/url"

	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/routers/api/repo"
	"github.com/code-devel-cover/CodeCover/routers/api/report"
	"github.com/code-devel-cover/CodeCover/routers/api/request"
	"github.com/code-devel-cover/CodeCover/routers/api/scm"
	"github.com/code-devel-cover/CodeCover/routers/api/user"
	"github.com/code-devel-cover/CodeCover/routers/docs"
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

type APIRouter struct {
	Config  *config.Config
	Session core.Session
	// service
	CoverageService core.CoverageService
	RepoService     core.RepoService
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

func (r *APIRouter) RegisterRoutes(e *gin.Engine) {
	docs.SwaggerInfo.Host = host(r.Config.Server.Addr)
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	g := e.Group("/api/v1")
	{
		g := g.Group("/user")
		g.POST("", user.HandleCreate())
	}
	{
		g := g.Group("/report")
		g.POST("/:id/:type", report.HandleUpload(
			r.CoverageService,
			r.ReportStore,
		))
	}
	{
		g := g.Group("/repo")
		g.Use(request.CheckLogin(r.Session))
		g.POST("/", repo.HandleCreate(r.RepoStore))
		g.GET("/:scm/:namespace/:name", repo.HandleGet(r.RepoStore))
		g.PATCH("/:scm/:namespace/:name/report", repo.HandleReportIDRenew(r.RepoStore, r.RepoService))
	}
	{
		g := g.Group("/scm")
		g.Use(request.CheckLogin(r.Session))
		g.GET("/:scm/repos", scm.HandleListSCM(r.RepoService))
	}
}
