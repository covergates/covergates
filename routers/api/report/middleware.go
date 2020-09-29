package report

import (
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/routers/api/request"
	"github.com/gin-gonic/gin"
)

// ProtectReport from modifying by unauthorized users
func ProtectReport(checkLogin gin.HandlerFunc, repoStore core.RepoStore, service core.SCMService) gin.HandlerFunc {
	return func(c *gin.Context) {
		setting := MustGetSetting(c)
		if !setting.Protected {
			return
		}
		checkLogin(c)
		if c.IsAborted() {
			return
		}
		ctx := c.Request.Context()
		user := request.MustGetUserFrom(c)
		repo := MustGetRepo(c)
		client, err := service.Client(repo.SCM)
		if err != nil {
			c.String(500, err.Error())
			c.Abort()
			return
		}
		creator, err := repoStore.Creator(repo)
		if err != nil {
			c.String(500, err.Error())
			c.Abort()
			return
		}
		if !client.Repositories().IsAdmin(ctx, user, repo.FullName()) && user.Login != creator.Login {
			c.String(401, "permission denied")
			c.Abort()
			return
		}
	}
}

// InjectReportContext such as repository, setting according to report id
func InjectReportContext(repoStore core.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		reportID := c.Param("id")
		repo, err := repoStore.Find(&core.Repo{ReportID: reportID})
		if err != nil {
			return
		}
		WithRepo(c, repo)
		if setting, err := repoStore.Setting(repo); err == nil {
			WithSetting(c, setting)
		}
	}
}
