package user

import (
	"strings"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/routers/api/request"
	"github.com/gin-gonic/gin"
)

// User for API response
type User struct {
	Login  string `json:"login"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

// Providers defines lined SCM accounts
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
			case core.GitLab:
				login = user.GitLabLogin
			case core.Bitbucket:
				login = user.BitbucketLogin
			default:
				login = ""
			}
			providers[strings.ToLower(string(scm))] = login != ""
		}
		c.JSON(200, providers)
	}
}

// HandleGetOwner of the repository
// @Summary Get repository's owner
// @Tags User
// @Param scm path string true "SCM"
// @Param namespace path string true "Namespace"
// @Param name path string true "name"
// @Success 200 {object} User "owner"
// @Router /user/owner/{scm}/{namespace}/{name} [get]
func HandleGetOwner(store core.RepoStore, service core.SCMService) gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := core.SCMProvider(c.Param("scm"))
		repo, err := store.Find(&core.Repo{
			NameSpace: c.Param("namespace"),
			Name:      c.Param("name"),
			SCM:       provider,
		})
		if err != nil {
			c.JSON(404, &User{})
			return
		}
		user, ok := request.UserFrom(c)
		if !ok {
			c.JSON(401, &User{})
			return
		}

		client, err := service.Client(provider)
		if err != nil {
			c.Error(err)
			c.JSON(500, &User{})
			return
		}

		if client.Repositories().IsAdmin(c.Request.Context(), user, repo.FullName()) {
			c.JSON(200, &User{
				Avatar: user.Avatar,
				Email:  user.Email,
				Login:  user.Login,
			})
			return
		}

		owner, err := store.Creator(repo)
		if err != nil {
			c.JSON(404, &User{})
			return
		}
		if owner.Login != user.Login {
			c.JSON(401, &User{})
			return
		}
		c.JSON(200, &User{
			Avatar: owner.Avatar,
			Email:  owner.Email,
			Login:  owner.Login,
		})
	}
}
