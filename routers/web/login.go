package web

import (
	"context"
	"net/http"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-login/login"
	"github.com/gin-gonic/gin"
)

const (
	keyLogin   = "login"
	keyAccess  = "access"
	keyRefresh = "refresh"
	keyExpires = "expires"
)

func createUser(ctx context.Context, service core.UserService, scm core.SCMProvider) (*core.User, error) {
	if err := service.Create(ctx, scm); err != nil {
		return nil, err
	}
	return service.Find(ctx, scm)
}

func HandleLogin(
	scm core.SCMProvider,
	userService core.UserService,
	session core.Session,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := WithToken(c)
		user, err := userService.Find(ctx, scm)
		if err != nil {
			user, err = createUser(ctx, userService, scm)
		}
		if err != nil {
			c.Error(err)
			c.String(400, err.Error())
		}
		if err := session.Create(c, user); err != nil {
			c.Error(err)
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
