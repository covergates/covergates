package report

import (
	"bytes"
	"fmt"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
func HandleUpload(
	scmService core.SCMService,
	coverageService core.CoverageService,
	repoStore core.RepoStore,
	reportStore core.ReportStore,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		reportID := c.Param("id")
		reportType := core.ReportType(c.Param("type"))
		commit, ok := c.GetPostForm("commit")
		if !ok {
			c.String(400, "must have commit SHA")
			return
		}
		ctx := c.Request.Context()

		repo, err := repoStore.Find(&core.Repo{ReportID: reportID})
		if err != nil {
			c.String(400, "cannot find repository related to report id")
			return
		}

		user, err := repoStore.Creator(repo)
		if err != nil {
			c.String(403, "repository creator not found")
		}

		client, err := scmService.Client(repo.SCM)
		if err != nil {
			c.String(400, err.Error())
			return
		}

		files, err := client.Contents().ListAllFiles(
			ctx, user,
			fmt.Sprintf("%s/%s", repo.NameSpace, repo.Name), commit)
		if err != nil {
			c.String(400, err.Error())
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.Error(err)
			c.String(400, err.Error())
			return
		}
		reader, err := file.Open()
		coverage, err := coverageService.Report(ctx, reportType, reader)
		if err != nil {
			c.String(400, err.Error())
			return
		}

		setting, err := repoStore.Setting(repo)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		if err := coverageService.TrimFileNames(ctx, coverage, setting.Filters); err != nil {
			c.String(500, err.Error())
		}

		report := &core.Report{
			ReportID: reportID,
			Coverage: coverage,
			Files:    files,
			Type:     reportType,
			Branch:   c.PostForm("branch"),
			Tag:      c.PostForm("tag"),
			Commit:   commit,
		}
		if err := reportStore.Upload(report); err != nil {
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

type getOptions struct {
	Latest bool `form:"latest"`
}

// HandleGet for the report id
// @Summary get reports for the report id
// @Tags Report
// @Param id path string true "report id"
// @Param latest query bool false "get latest report in main branch"
// @Success 200 {object} core.Report "coverage report"
// @Router /reports/{id} [get]
func HandleGet(reportStore core.ReportStore, repoStore core.RepoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		reportID := c.Param("id")
		option := &getOptions{}
		if err := c.BindQuery(option); err != nil {
			log.Error(err)
			c.JSON(400, []*core.Report{})
			return
		}
		// TODO: support multiple type (language) reports in one repository
		if option.Latest {
			report, err := getLatest(reportStore, repoStore, reportID)
			if err != nil {
				c.JSON(404, []*core.Report{})
			} else {
				c.JSON(200, []*core.Report{report})
			}
		} else {
			reports, err := getAll(reportStore, reportID)
			if err != nil {
				c.JSON(400, []*core.Report{})
			} else {
				c.JSON(200, reports)
			}
		}
		return
	}
}

// HandleGetTreeMap for coverage difference with main branch
// @Summary Get coverage difference treemap with main branch
// @Tags Report
// @Produce image/svg+xml
// @Param id path string true "report id"
// @param commit path string true "commit sha"
// @Success 200 {object} string "treemap svg"
// @Router /reports/{id}/{commit}/treemap [get]
func HandleGetTreeMap(
	reportStore core.ReportStore,
	repoStore core.RepoStore,
	chartService core.ChartService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		reportID := c.Param("id")
		commit := c.Param("commit")
		new, err := reportStore.Find(&core.Report{
			ReportID: reportID,
			Commit:   commit,
		})
		if err != nil {
			c.String(500, err.Error())
			return
		}
		old, err := getLatest(reportStore, repoStore, reportID)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		chart := chartService.CoverageDiffTreeMap(old.Coverage, new.Coverage)
		buffer := bytes.NewBuffer([]byte{})
		if err := chart.Render(buffer); err != nil {
			c.String(500, err.Error())
			return
		}
		c.Data(200, "image/svg+xml", buffer.Bytes())
		return
	}
}

func getLatest(reportStore core.ReportStore, repoStore core.RepoStore, reportID string) (*core.Report, error) {
	repo, err := repoStore.Find(&core.Repo{ReportID: reportID})
	if err != nil {
		return nil, err
	}
	return reportStore.Find(&core.Report{
		ReportID: reportID,
		Branch:   repo.Branch,
	})
}

func getAll(store core.ReportStore, reportID string) ([]*core.Report, error) {
	return store.Finds(&core.Report{
		ReportID: reportID,
	})
}
