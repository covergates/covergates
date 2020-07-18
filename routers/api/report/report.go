package report

import (
	"bytes"
	"context"
	"io"

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
		branch := c.PostForm("branch")
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

		setting, err := repoStore.Setting(repo)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		// get upload file
		file, err := c.FormFile("file")
		if err != nil {
			c.Error(err)
			c.String(400, err.Error())
			return
		}

		gitRepo, err := getGitRepository(ctx, repoStore, scmService, repo)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		files, err := gitRepo.ListAllFiles(commit)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		if branch == "" {
			co, err := gitRepo.Commit(commit)
			if err == nil && co.InDefaultBranch() {
				branch = repo.Branch
			}
		}

		reader, err := file.Open()
		coverage, err := loadCoverageReport(ctx, coverageService, reportType, reader, setting)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		report := &core.Report{
			ReportID: reportID,
			Coverage: coverage,
			Files:    files,
			Type:     reportType,
			Branch:   branch,
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

// getAll reports related to gitven reportID
func getAll(store core.ReportStore, reportID string) ([]*core.Report, error) {
	return store.Finds(&core.Report{
		ReportID: reportID,
	})
}

// loadCoverageReort from io reader and apply repository wide setting
func loadCoverageReport(
	ctx context.Context,
	service core.CoverageService,
	reportType core.ReportType,
	data io.Reader,
	setting *core.RepoSetting,
) (*core.CoverageReport, error) {
	coverage, err := service.Report(ctx, reportType, data)
	if err != nil {
		return nil, err
	}
	if err := service.TrimFileNames(ctx, coverage, setting.Filters); err != nil {
		return nil, err
	}
	return coverage, nil
}

// getGitRepository with given Repo
func getGitRepository(
	ctx context.Context,
	store core.RepoStore,
	service core.SCMService,
	repo *core.Repo,
) (core.GitRepository, error) {
	user, err := store.Creator(repo)
	if err != nil {
		return nil, err
	}
	client, err := service.Client(repo.SCM)
	if err != nil {
		return nil, err
	}
	return client.Git().GitRepository(ctx, user, repo.FullName())
}
