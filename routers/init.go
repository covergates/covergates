package routers

import (
	"github.com/code-devel-cover/CodeCover/routers/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.Use(cors.Default())
	web.RegisterStaticWebHandlers(r)
	r.NoRoute(web.HandleIndex)
}
