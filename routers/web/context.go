package web

import (
	"context"

	"github.com/drone/go-scm/scm"
	"github.com/gin-gonic/gin"
)

func WithToken(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	ctx = scm.WithContext(ctx, &scm.Token{
		Token:   c.GetString(keyAccess),
		Refresh: c.GetString(keyRefresh),
		Expires: c.GetTime(keyExpires),
	})
	return ctx
}
