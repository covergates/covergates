package report

import (
	"github.com/covergates/covergates/core"
	"github.com/gin-gonic/gin"
)

const (
	keyRepo    = "report_repo"
	keySetting = "report_setting"
)

// WithRepo context
func WithRepo(c *gin.Context, repo *core.Repo) {
	c.Set(keyRepo, repo)
}

// MustGetRepo from context
func MustGetRepo(c *gin.Context) *core.Repo {
	return c.MustGet(keyRepo).(*core.Repo)
}

// WithSetting context
func WithSetting(c *gin.Context, setting *core.RepoSetting) {
	c.Set(keySetting, setting)
}

// MustGetSetting from context
func MustGetSetting(c *gin.Context) *core.RepoSetting {
	return c.MustGet(keySetting).(*core.RepoSetting)
}
