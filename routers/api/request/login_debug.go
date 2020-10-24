//+build debug

package request

import (
	"os"

	"github.com/covergates/covergates/core"
	"github.com/gin-gonic/gin"
)

func CheckLogin(session core.Session, oauth core.OAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &core.User{
			Login:          os.Getenv("DEBUG_LOGIN"),
			Email:          os.Getenv("DEBUG_EMAIL"),
			Avatar:         os.Getenv("DEBUG_AVATAR"),
			GithubLogin:    os.Getenv("DEBUG_GITHUB_LOGIN"),
			GithubEmail:    os.Getenv("DEBUG_EMAIL"),
			GithubToken:    os.Getenv("DEBUG_GITHUB_TOKEN"),
			GiteaLogin:     os.Getenv("DEBUG_GITEA_LOGIN"),
			GiteaEmail:     os.Getenv("DEBUG_EMAIL"),
			GiteaToken:     os.Getenv("DEBUG_GITEA_TOKEN"),
			GitLabLogin:    os.Getenv("DEBUG_GITLAB_LOGIN"),
			GitLabEmail:    os.Getenv("DEBUG_EMAIL"),
			GitLabToken:    os.Getenv("DEBUG_GITLAB_TOKEN"),
			BitbucketLogin: os.Getenv("DEBUG_BITBUCKET_LOGIN"),
			BitbucketEmail: os.Getenv("DEBUG_EMAIL"),
			BitbucketToken: os.Getenv("DEBUG_BITBUCKET_TOKEN"),
		}
		WithUser(c, user)
	}
}
