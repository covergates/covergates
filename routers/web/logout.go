package web

import (
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/gin-gonic/gin"
)

func RegisterLogout(r *gin.Engine, session core.Session) {
	r.Any("/logout", handleLogout(session))
}

func handleLogout(session core.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := session.Clear(c); err != nil {
			c.Error(err)
			c.String(500, "Fail to logout")
			return
		}
	}
}
