package web

import (
	"net/http"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-login/login"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	keyLogin   = "login"
	keyAccess  = "access"
	keyRefresh = "refresh"
	keyExpires = "expires"
)

func HandleLogin(
	scm core.SCMProvider,
	scmService core.SCMService,
	session core.Session,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.GetBool(keyLogin) {
			return
		}
		client, err := scmService.Client(scm)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		ctx := c.Request.Context()
		user, err := client.Users().Find(ctx, TokenFrom(c))
		if err != nil {
			user, err = client.Users().Create(ctx, TokenFrom(c))
		}
		if err != nil {
			log.Error(err)
			c.String(400, err.Error())
		}
		if err := session.Create(c, user); err != nil {
			log.Error(err)
			c.String(400, err.Error())
			return
		}
		c.String(200, "login")
	}
}

func MiddlewareLogin(scm core.SCMProvider, m core.LoginMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := m.Handler(scm)
		h := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			err := login.ErrorFrom(ctx)
			if err != nil {
				c.Error(err)
				c.Abort()
				return
			}
			tok := login.TokenFrom(ctx)
			c.Set(keyLogin, true)
			c.Set(keyAccess, tok.Access)
			c.Set(keyExpires, tok.Expires)
			c.Set(keyRefresh, tok.Refresh)
		}))
		h.ServeHTTP(c.Writer, c.Request)
	}
}
