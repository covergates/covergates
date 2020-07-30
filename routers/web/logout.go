package web

import (
	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/gin-gonic/gin"
)

// HandleLogout user session
func HandleLogout(config *config.Config, session core.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := session.Clear(c); err != nil {
			c.Error(err)
			c.String(500, "Fail to logout")
			return
		}
		c.Header("Cache-Control", "no-store")
		c.Redirect(301, config.Server.BaseURL())
	}
}
