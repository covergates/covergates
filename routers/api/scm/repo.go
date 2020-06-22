package scm

import (
	"errors"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/routers/api/request"
	"github.com/gin-gonic/gin"
)

var errUserNotFound = errors.New("user not found")

// HandleListSCM request, returns repositories
// @Summary Get repositories from SCM
// @Tags SCM
// @Param scm path string true "SCM source (github, gitea)"
// @Success 200 {object} []core.Repo "repositories"
// @Router /scm/{scm}/repos [get]
func HandleListSCM(service core.SCMService) gin.HandlerFunc {
	return func(c *gin.Context) {
		scm := core.SCMProvider(c.Param("scm"))
		user, ok := request.UserFrom(c)
		if !ok {
			c.JSON(500, errUserNotFound.Error())
			return
		}
		ctx := c.Request.Context()
		client, err := service.Client(scm)
		if err != nil {
			c.Error(err)
			c.JSON(500, []*core.Repo{})
			return
		}
		repositories, err := client.Repositories().List(ctx, user)
		if err != nil {
			c.Error(err)
			c.JSON(500, []*core.Repo{})
			return
		}
		c.JSON(200, repositories)
	}
}
