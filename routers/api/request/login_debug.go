//+build debug

package request

import (
	"os"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/gin-gonic/gin"
)

func CheckLogin(session core.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &core.User{
			Login:       os.Getenv("DEBUG_LOGIN"),
			Email:       os.Getenv("DEBUG_EMAIL"),
			GithubLogin: os.Getenv("DEBUG_LOGIN"),
			GithubEmail: os.Getenv("DEBUG_EMAIL"),
			GithubToken: os.Getenv("DEBUG_GITHUB_TOKEN"),
			GiteaEmail:  os.Getenv("DEBUG_EMAIL"),
			GiteaToken:  os.Getenv("DEBUG_GITEA_TOKEN"),
		}
		WithUser(c, user)
	}
}
