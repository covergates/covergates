package routers

import (
	"github.com/code-devel-cover/CodeCover/core"
	_ "github.com/code-devel-cover/CodeCover/models"
	"github.com/code-devel-cover/CodeCover/routers/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:generate swag init -g ./init.go -d ./ -o ./docs

// @title CodeCover API
// @version 1.0
// @description REST API for CodeCover
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

type Routers struct {
	LoginMiddleware core.LoginMiddleware
}

func (router *Routers) RegisterRoutes(r *gin.Engine) {
	r.Use(cors.Default())
	web.RegisterStaticWeb(r)
	web.RegisterLogin(r, router.LoginMiddleware)
	r.NoRoute(web.HandleIndex)
}
