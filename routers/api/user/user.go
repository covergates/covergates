package user

import (
	"strings"

	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/routers/api/request"
	"github.com/gin-gonic/gin"
)

// User for API response
type User struct {
	Login  string `json:"login"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

type Providers map[string]bool

// HandleCreate new user
func HandleCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement user create
	}
}

// HandleGet login user
// @Summary Get login user
// @Tags User
// @Success 200 {object} User "user"
// @Failure 404 {string} string "error"
// @Router /user [get]
func HandleGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := request.UserFrom(c)
		if !ok {
			c.String(404, "user not found")
			return
		}
		u := User{
			Login:  user.Login,
			Email:  user.Email,
			Avatar: user.Avatar,
		}
		c.JSON(200, u)
	}
}

// HandleGetSCM binds to login user
// @Summary Get user's SCM binding state
// @Tags User
// @Success 200 {object} Providers "providers"
// @Failure 404 {object} Providers "providers"
// @Router /user/scm [get]
func HandleGetSCM(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		providers := make(Providers)
		user, ok := request.UserFrom(c)
		if !ok {
			c.JSON(404, providers)
			return
		}
		for _, scm := range config.Providers() {
			login := ""
			switch scm {
			case core.Gitea:
				login = user.GiteaLogin
			case core.Github:
				login = user.GithubLogin
			default:
				login = ""
			}
			providers[strings.ToLower(string(scm))] = login != ""
		}
		c.JSON(200, providers)
	}
}
