package repo

import (
	"fmt"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/routers/api/request"
	"github.com/gin-gonic/gin"
)

// TODO: Need to check user permission

// HandleCreate a repository
// @Summary Create new repository for code coverage
// @Tags Repository
// @Param repo body core.Repo true "repository to create"
// @Success 200 {object} core.Repo "Created repository"
// @Router /repos [post]
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
// @Router /repos/{scm}/{namespace}/{name}/report [patch]
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
// @Router /repos/{scm}/{namespace}/{name} [get]
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

// HandleListSCM repositories
// @Summary List repositories
// @Tags Repository
// @Param scm path string true "SCM source (github, gitea)"
// @Success 200 {object} []core.Repo "repositories"
// @Router /repos/{scm} [get]
func HandleListSCM(service core.SCMService, store core.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		scm := core.SCMProvider(c.Param("scm"))
		user, ok := request.UserFrom(c)
		if !ok {
			c.JSON(401, []*core.Repo{})
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
		urls := make([]string, len(repositories))
		for i, repo := range repositories {
			urls[i] = repo.URL
		}
		storeRepositories, err := store.Finds(urls...)
		if err != nil {
			c.JSON(200, repositories)
			return
		}
		reportsMap := make(map[string]string)
		for _, repo := range storeRepositories {
			reportsMap[repo.URL] = repo.ReportID
		}
		for _, repo := range repositories {
			report, ok := reportsMap[repo.URL]
			if ok {
				repo.ReportID = report
			}
		}
		c.JSON(200, repositories)
	}
}

// HandleGetFiles from the repository
// @Summary List all files in repository
// @Tags Repository
// @Param scm path string true "SCM"
// @Param namespace path string true "Namespace"
// @Param name path string true "name"
// @Param ref query string false "specify git ref, default main branch"
// @Success 200 {object} []string "files"
// @Router /repos/{scm}/{namespace}/{name}/files [get]
func HandleGetFiles(service core.SCMService) gin.HandlerFunc {
	return func(c *gin.Context) {
		scm := core.SCMProvider(c.Param("scm"))
		repoName := fmt.Sprintf("%s/%s", c.Param("namespace"), c.Param("name"))
		user, ok := request.UserFrom(c)
		ctx := c.Request.Context()
		if !ok {
			c.JSON(401, []string{})
			return
		}
		client, err := service.Client(scm)
		if err != nil {
			c.Error(err)
			c.JSON(500, []string{})
			return
		}
		ref := c.Query("ref")
		if ref == "" {
			repo, err := client.Repositories().Find(ctx, user, repoName)
			if err != nil {
				c.Error(err)
				c.JSON(500, []string{})
				return
			}
			ref = repo.Branch
		}
		files, err := client.Contents().ListAllFiles(ctx, user, repoName, ref)
		if err != nil {
			c.Error(err)
			c.JSON(500, []string{})
			return
		}
		c.JSON(200, files)
	}
}
