package repo

import (
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/gin-gonic/gin"
)

// TODO: Need to check user permission

// HandleCreate a repository
// @Summary Create new repository for code coverage
// @Tags Repository
// @Param repo body core.Repo true "repository to create"
// @Success 200 {object} core.Repo "Created repository"
// @Router /repo [post]
func HandleCreate(store core.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := &core.Repo{}
		if err := c.BindJSON(repo); err != nil {
			c.String(400, err.Error())
			return
		}
		if err := store.Create(repo); err != nil {
			c.String(400, err.Error())
			return
		}
		c.JSON(200, repo)
	}
}

// HandleReportIDRenew generates a new report id for the repository
// @Summary renew repository report id
// @Tags Repository
// @Param scm path string true "SCM"
// @Param namespace path string true "Namespace"
// @Param name path string true "name"
// @Success 200 {object} core.Repo "updated repository"
// @Router /repo/{scm}/{namespace}/{name}/report [patch]
func HandleReportIDRenew(store core.RepoStore, service core.SCMService) gin.HandlerFunc {
	return func(c *gin.Context) {
		scm := core.SCMProvider(c.Param("scm"))
		repo, err := store.Find(&core.Repo{
			Name:      c.Param("name"),
			NameSpace: c.Param("namespace"),
			SCM:       scm,
		})
		if err != nil {
			c.String(500, err.Error())
			return
		}
		client, err := service.Client(scm)
		if err != nil {
			c.String(500, err.Error())
		}
		repo.ReportID = client.Repositories().NewReportID(repo)
		if err := store.Update(repo); err != nil {
			c.String(500, err.Error())
			return
		}
		c.JSON(200, repo)
		return
	}
}

// HandleGet repository
// @Summary get repository
// @Tags Repository
// @Param scm path string true "SCM"
// @Param namespace path string true "Namespace"
// @Param name path string true "name"
// @Success 200 {object} core.Repo found repository
// @Router /repo/{scm}/{namespace}/{name} [get]
func HandleGet(store core.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo, err := store.Find(&core.Repo{
			NameSpace: c.Param("namespace"),
			Name:      c.Param("name"),
			SCM:       core.SCMProvider(c.Param("scm")),
		})
		if err != nil {
			c.JSON(404, &core.Repo{})
			return
		}
		c.JSON(200, repo)
	}
}
