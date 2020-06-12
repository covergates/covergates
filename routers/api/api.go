package api

import (
	"github.com/code-devel-cover/CodeCover/routers/api/user"
	"github.com/gin-gonic/gin"
)

type APIRouter struct {
}

func (r *APIRouter) RegisterRoutes(e *gin.Engine) {
	g := e.Group("/api/v1")
	{
		g := g.Group("/user")
		g.POST("", user.HandleCreate())
	}
}
