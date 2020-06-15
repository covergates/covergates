package report

import (
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/gin-gonic/gin"
)

// HandleUpload report
// @Summary Upload coverage report
// @Tags Report
// @Param type path string true "report type"
// @Param id path string	true "report id"
// @Param file formData file true "report"
// @Param commit formData string true "Git commit SHA"
// @Param branch formData string false "branch ref"
// @Param tag formData string false "tag ref"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "error message"
// @Router /reports/{id}/{type} [post]
func HandleUpload(service core.CoverageService, store core.ReportStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		reportID := c.Param("id")
		reportType := core.ReportType(c.Param("type"))
		commit, ok := c.GetPostForm("commit")
		if !ok {
			c.String(400, "must have commit SHA")
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			c.Error(err)
			c.String(400, err.Error())
			return
		}
		ctx := c.Request.Context()
		reader, err := file.Open()
		coverage, err := service.Report(ctx, reportType, reader)
		if err != nil {
			c.String(400, err.Error())
			return
		}
		report := &core.Report{
			ReportID: reportID,
			Coverage: coverage,
			Type:     reportType,
			Branch:   c.PostForm("branch"),
			Tag:      c.PostForm("tag"),
			Commit:   commit,
		}
		if err := store.Upload(report); err != nil {
			c.Error(err)
			c.String(500, err.Error())
			return
		}
		c.String(200, "ok")
	}
}

// HandleRepo for report id
// @Summary get repository of the report id
// @Tags Report
// @Param id path string true "report id"
// @Success 200 {object} core.Repo "repository"
// @Failure 400 {string} string "error message"
// @Router /reports/{id}/repo [get]
func HandleRepo(store core.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo, err := store.Find(&core.Repo{
			ReportID: c.Param("id"),
		})
		if err != nil {
			c.JSON(404, &core.Repo{})
			return
		}
		c.JSON(200, repo)
	}
}

// HandleGet for the report id
// @Summary get reports for the report id
// @Tags Report
// @Param id path string true "report id"
// @Param latest query bool false "get latest report"
// @Success 200 {object} core.Report "coverage report"
// @Router /reports/{id} [get]
func HandleGet(store core.ReportStore) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func getLatest(store core.ReportStore, reportID string) (*core.Report, error) {
	// store.Find(reportID
	return nil, nil
}

// TODO:implements get all reports
func getAll(store core.ReportStore) ([]*core.Report, error) {
	return nil, nil
}
