package repo

import (
	"fmt"
	"strings"

	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/routers/api/request"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

// HandleCreate a repository
// @Summary Create new repository for code coverage
// @Tags Repository
// @Param repo body core.Repo true "repository to create"
// @Success 200 {object} core.Repo "Created repository"
// @Router /repos [post]
func HandleCreate(store core.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := request.UserFrom(c)
		if !ok {
			c.String(403, "user not found")
			return
		}
		repo := &core.Repo{}
		if err := c.BindJSON(repo); err != nil {
			c.String(400, err.Error())
			return
		}
		if err := store.Create(repo, user); err != nil {
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

// HandleListAll repositories
// @Summary List repositories for all available SCM providers
// @Tags Repository
// @Success 200 {object} []core.Repo "repositories"
// @Router /repos [get]
func HandleListAll(config *config.Config, service core.SCMService, store core.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		repositories := make([]*core.Repo, 0)
		user, ok := request.UserFrom(c)
		if !ok {
			c.JSON(401, repositories)
			return
		}
		ctx := c.Request.Context()
		for _, provider := range config.Providers() {
			client, err := service.Client(provider)
			if err != nil {
				continue
			}
			result, err := getRepositories(ctx, user, client, store)
			if err != nil {
				continue
			}
			repositories = append(repositories, result...)
		}
		c.JSON(200, repositories)
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
		repositories, err := getRepositories(ctx, user, client, store)
		if err != nil {
			c.JSON(500, []*core.Repo{})
			return
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
		ref, err := getRef(c, client, user)
		if err != nil {
			c.Error(err)
			c.JSON(500, []string{})
			return
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

// HandleGetFileContent for a file from repository
// @Summary Get a file content
// @Tags Repository
// @Param scm path string true "SCM"
// @Param namespace path string true "Namespace"
// @Param name path string true "name"
// @Param path path string true "file path"
// @Param ref query string false "specify git ref, default main branch"
// @Success 200 {string} string "file content"
// @Router /repos/{scm}/{namespace}/{name}/content/{path} [get]
func HandleGetFileContent(service core.SCMService) gin.HandlerFunc {
	return func(c *gin.Context) {
		scm := core.SCMProvider(c.Param("scm"))
		repoName := fmt.Sprintf("%s/%s", c.Param("namespace"), c.Param("name"))
		filePath := strings.TrimLeft(c.Param("path"), "/")
		user, ok := request.UserFrom(c)
		ctx := c.Request.Context()
		if !ok {
			c.String(401, "")
			return
		}
		client, err := service.Client(scm)
		if err != nil {
			c.Error(err)
			c.String(500, "")
			return
		}
		ref, err := getRef(c, client, user)
		if err != nil {
			c.Error(err)
			c.String(500, "")
			return
		}
		content, err := client.Contents().Find(ctx, user, repoName, filePath, ref)
		c.Data(200, "text/plain", content)
	}
}

// HandleGetSetting for the repository
// @Summary get repository setting
// @Tags Repository
// @Param scm path string true "SCM"
// @Param namespace path string true "Namespace"
// @Param name path string true "name"
// @Success 200 {object} core.RepoSetting repository setting
// @Router /repos/{scm}/{namespace}/{name}/setting [get]
func HandleGetSetting(store core.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo, err := store.Find(&core.Repo{
			NameSpace: c.Param("namespace"),
			Name:      c.Param("name"),
			SCM:       core.SCMProvider(c.Param("scm")),
		})
		if err != nil {
			c.JSON(404, &core.RepoSetting{})
			return
		}
		setting, err := store.Setting(repo)
		if err != nil {
			if err != nil {
				c.JSON(404, &core.RepoSetting{})
				return
			}
		}
		c.JSON(200, setting)
	}
}

// HandleUpdateSetting for the repository
// @Summary get repository setting
// @Tags Repository
// @Param scm path string true "SCM"
// @Param namespace path string true "Namespace"
// @Param name path string true "name"
// @Param setting body core.RepoSetting true "repository setting"
// @Success 200 {object} core.RepoSetting repository setting
// @Router /repos/{scm}/{namespace}/{name}/setting [post]
func HandleUpdateSetting(store core.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo, err := store.Find(&core.Repo{
			NameSpace: c.Param("namespace"),
			Name:      c.Param("name"),
			SCM:       core.SCMProvider(c.Param("scm")),
		})
		if err != nil {
			c.JSON(404, &core.RepoSetting{})
			return
		}
		setting := &core.RepoSetting{}
		if err := c.BindJSON(setting); err != nil {
			c.JSON(400, setting)
			return
		}
		if err := store.UpdateSetting(repo, setting); err != nil {
			c.JSON(500, setting)
			return
		}
		c.JSON(200, setting)
	}
}

func getRef(c *gin.Context, client core.Client, user *core.User) (string, error) {
	repoName := fmt.Sprintf("%s/%s", c.Param("namespace"), c.Param("name"))
	ref := c.Query("ref")
	if ref == "" {
		repo, err := client.Repositories().Find(c.Request.Context(), user, repoName)
		if err != nil {
			return "", err
		}
		ref = repo.Branch
	}
	return ref, nil
}

func getRepositories(
	ctx context.Context,
	user *core.User,
	client core.Client,
	store core.RepoStore,
) ([]*core.Repo, error) {
	repositories, err := client.Repositories().List(ctx, user)
	if err != nil {
		return nil, err
	}
	urls := make([]string, len(repositories))
	for i, repo := range repositories {
		urls[i] = repo.URL
	}
	storeRepositories, err := store.Finds(urls...)
	if err != nil {
		return repositories, nil
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
	return repositories, nil
}
