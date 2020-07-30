package web

import (
	"github.com/covergates/covergates/core"
	"github.com/gin-gonic/gin"
)

// TokenFrom context
func TokenFrom(c *gin.Context) *core.Token {
	return &core.Token{
		Token:   c.GetString(keyAccess),
		Refresh: c.GetString(keyRefresh),
		Expires: c.GetTime(keyExpires),
	}
}
