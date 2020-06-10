package web

import (
	"log"
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

func RegisterLogin(
	r *gin.Engine,
	m core.LoginMiddleware,
) {
	{
		g := r.Group("/login")
		g.Any("/github", MiddlewareLogin(core.Github, m), HandleLogin)
	}
}

func HandleLogin(c *gin.Context) {
	if !c.GetBool("login") {
		log.Println("Not login")
		return
	}
	log.Printf("%s", c.GetString(keyAccess))
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
