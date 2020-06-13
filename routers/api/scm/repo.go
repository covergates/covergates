package scm

import (
	"errors"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/routers/api/request"
	"github.com/gin-gonic/gin"
)

var errUserNotFound = errors.New("user not found")

type Repo struct{}

// HandleList request, returns repositories
// @Summary Get repositories from SCM
// @Param scm path string true "SCM source (github, gitea)"
// @Success 200 {object} []core.Repo "repositories"
// @Router /scm/{scm}/repos [get]
func HandleList(service core.RepoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		scm := core.SCMProvider(c.Param("scm"))
		user, ok := request.UserFrom(c)
		if !ok {
			c.JSON(500, errUserNotFound)
			return
		}
		ctx := c.Request.Context()
		repositories, err := service.List(ctx, scm, user)
		if err != nil {
			c.JSON(500, err)
			return
		}
		c.JSON(200, repositories)
	}
}
