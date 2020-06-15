package web

import (
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/gin-gonic/gin"
)

func TokenFrom(c *gin.Context) *core.Token {
	return &core.Token{
		Token:   c.GetString(keyAccess),
		Refresh: c.GetString(keyRefresh),
		Expires: c.GetTime(keyExpires),
	}
}
