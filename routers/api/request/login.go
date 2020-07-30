//+build !debug

package request

import (
	"github.com/covergates/covergates/core"
	"github.com/gin-gonic/gin"
)

// CheckLogin session
func CheckLogin(session core.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := session.GetUser(c)
		if user.Login == "" {
			c.String(401, "Unauthorized")
			c.Abort()
			return
		}
		WithUser(c, user)
	}
}
