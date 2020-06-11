package web

import (
	"errors"
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

var errClientNotFound = errors.New("SCM Client not found")

func RegisterLogin(
	r *gin.Engine,
	middleware core.LoginMiddleware,
	clientService core.SCMClientService,
	userService core.UserService,
	session core.Session,
) {
	{
		g := r.Group("/login")
		g.Any("/github",
			MiddlewareLogin(core.Github, middleware),
			HandleLogin(
				core.Github,
				clientService,
				userService,
				session,
			),
		)
	}
}

func HandleLogin(
	scm core.SCMProvider,
	clientService core.SCMClientService,
	userService core.UserService,
	session core.Session,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := clientService.Client(scm)
		if client == nil {
			c.Error(errClientNotFound)
			c.Abort()
			return
		}
		ctx := WithToken(c)
		user, err := userService.Find(ctx, scm)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		if err := session.Create(c, user); err != nil {
			c.Error(err)
			c.Abort()
			return
		}
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
