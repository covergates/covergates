package request

import (
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/gin-gonic/gin"
)

func CheckLogin(session core.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := session.Get(c)
		if user.Login == "" {
			c.String(401, "Unauthorized")
			c.Abort()
			return
		}
		WithUser(c, user)
	}
}
