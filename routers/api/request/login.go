package request

import (
	"github.com/code-devel-cover/CodeCover/modules/session"
	"github.com/gin-gonic/gin"
)

func CheckLogin(session *session.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := session.Get(c)
		if user.Login == "" {
			c.String(403, "Unauthorized")
			c.Abort()
			return
		}
		WithUser(c, user)
	}
}
