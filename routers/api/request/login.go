//+build !debug

package request

import (
	"github.com/covergates/covergates/core"
	"github.com/gin-gonic/gin"
)

// CheckLogin session
func CheckLogin(session core.Session, oauth core.OAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := oauth.Validate(c.Request)
		if err == nil && user != nil && user.Login != "" {
			WithUser(c, user)
			return
		}
		user = session.GetUser(c)
		if user.Login == "" {
			c.String(401, "Unauthorized")
			c.Abort()
			return
		}
		WithUser(c, user)
	}
}
