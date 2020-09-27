package user

import (
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/routers/api/request"
	"github.com/gin-gonic/gin"
)

// HandleSynchronizeRepo for the user
// @Summary Synchronize user's repository from remote SCM
// @Tags User
// @Success 200 {string} string "ok"
// @Router /user/repos [patch]
func HandleSynchronizeRepo(repoService core.RepoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := request.MustGetUserFrom(c)
		ctx := c.Request.Context()
		if err := repoService.Synchronize(ctx, user); err != nil {
			c.Error(err)
			c.String(500, err.Error())
			return
		}
		c.String(200, "ok")
	}
}

// HandleListRepo for the uer
// @Summary List user synchronized repositories
// @Tags User
// @Success 200 {object} []core.Repo "list of repositories"
// @Router /user/repos [get]
func HandleListRepo(userStore core.UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := request.MustGetUserFrom(c)
		repos, err := userStore.ListRepositories(user)
		if err != nil {
			c.Error(err)
			c.JSON(500, []*core.Repo{})
			return
		}
		c.JSON(200, repos)
	}
}
