package api

import (
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/routers/api/report"
	"github.com/code-devel-cover/CodeCover/routers/api/user"
	_ "github.com/code-devel-cover/CodeCover/routers/docs"
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
	CoverageService core.CoverageService
	ReportStore     core.ReportStore
}

func (r *APIRouter) RegisterRoutes(e *gin.Engine) {
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
}
