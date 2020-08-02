package report

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/routers/api/request"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// HandleUpload report
// @Summary Upload coverage report
// @Tags Report
// @Param id path string	true "report id"
// @Param file formData file true "report"
// @Param commit formData string true "Git commit SHA"
// @Param type formData string true "report type"
// @Param branch formData string false "branch ref"
// @Param tag formData string false "tag ref"
// @Param files formData string false "files list of the repository"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "error message"
// @Router /reports/{id} [post]
func HandleUpload(
	coverageService core.CoverageService,
	repoStore core.RepoStore,
	reportStore core.ReportStore,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: need to handle repository update action according to setting
		if _, ok := c.GetPostForm("type"); !ok {
			c.String(400, "must have report type")
			return
		}

		if _, ok := c.GetPostForm("commit"); !ok {
			c.String(400, "must have commit SHA")
			return
		}

		reportID := c.Param("id")
		branch := c.PostForm("branch")
		reportType := core.ReportType(c.PostForm("type"))
		commit := c.PostForm("commit")
		ctx := c.Request.Context()

		repo, err := repoStore.Find(&core.Repo{ReportID: reportID})
		if err != nil {
			c.String(400, "cannot find repository related to report id")
			return
		}

		setting, err := repoStore.Setting(repo)
		if err != nil {
			log.Error(err)
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

		files := make([]string, 0)
		if c.PostForm("files") != "" {
			if err := json.Unmarshal([]byte(c.PostForm("files")), &files); err != nil {
				c.Error(err)
				c.String(500, err.Error())
				return
			}
		}

		reader, err := file.Open()
		coverage, err := loadCoverageReport(ctx, coverageService, reportType, reader, setting)
		if err != nil {
			log.Error(err)
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
func HandleGet(
	reportStore core.ReportStore,
	repoStore core.RepoStore,
	service core.SCMService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		reportID := c.Param("id")
		option := &getOptions{}
		if err := c.BindQuery(option); err != nil {
			log.Error(err)
			c.JSON(400, []*core.Report{})
			return
		}
		if !hasPermission(c, repoStore, service, reportID) {
			c.JSON(401, []*core.Report{})
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
// @param source path string true "source branch"
// @Success 200 {object} string "treemap svg"
// @Router /reports/{id}/treemap/{source} [get]
func HandleGetTreeMap(
	reportStore core.ReportStore,
	repoStore core.RepoStore,
	chartService core.ChartService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		reportID := c.Param("id")
		source := c.Param("source")
		new, err := reportStore.Find(&core.Report{
			ReportID: reportID,
			Branch:   source,
		})
		if err != nil {
			c.String(500, err.Error())
			return
		}
		old, err := getLatest(reportStore, repoStore, reportID)
		if err != nil {
			old = &core.Report{
				Coverage: &core.CoverageReport{},
			}
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

// HandleComment report summary
// @Summary Leave a report summary comment on pull request
// @Tags Report
// @Param id path string true "report id"
// @param number path string true "pull request number"
// @Success 200 {object} string "ok"
// @Router /reports/{id}/comment/{number} [POST]
func HandleComment(
	config *config.Config,
	service core.SCMService,
	repoStore core.RepoStore,
	reportStore core.ReportStore,
	reportService core.ReportService,
) gin.HandlerFunc {
	// TODO: Add handle comment unit test
	// TODO: Need to test comment with SHA or branch
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		reportID := c.Param("id")
		number, err := strconv.Atoi(c.Param("number"))
		if err != nil {
			c.String(400, "invalid pull request number")
			return
		}
		repo, err := repoStore.Find(&core.Repo{ReportID: reportID})
		if err != nil {
			c.String(400, "repository not found")
			return
		}
		user, err := repoStore.Creator(repo)
		if err != nil {
			c.String(400, "user not found")
			return
		}
		client, err := service.Client(repo.SCM)
		pr, err := client.PullRequests().Find(ctx, user, repo.FullName(), number)
		if err != nil {
			c.String(400, "cannot find pull request")
			return
		}

		// TODO: handle multiple language repository
		source, err := reportStore.Find(&core.Report{ReportID: reportID, Commit: pr.Commit})
		if err != nil {
			if source, err = reportStore.Find(&core.Report{ReportID: reportID, Branch: pr.Source}); err != nil {
				c.String(500, err.Error())
				return
			}
		}
		target, err := reportStore.Find(&core.Report{ReportID: reportID, Branch: pr.Target})
		if err != nil {
			target = &core.Report{}
		}

		r, err := reportService.MarkdownReport(source, target)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		buf := &bytes.Buffer{}

		buf.WriteString(fmt.Sprintf(
			"![treemap](%s/api/v1/reports/%s/treemap/%s?base=%s)\n\n",
			config.Server.URL(),
			reportID,
			source.Branch,
			target.Branch,
		),
		)

		if _, err := io.Copy(buf, r); err != nil {
			c.String(500, err.Error())
			return
		}

		if comment, err := reportStore.FindComment(&core.Report{ReportID: reportID}, number); err == nil {
			client.PullRequests().RemoveComment(ctx, user, repo.FullName(), number, comment.Comment)
		}

		commentID, err := client.PullRequests().CreateComment(
			ctx,
			user,
			repo.FullName(),
			number,
			string(buf.Bytes()),
		)
		log.Println(commentID)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		comment := &core.ReportComment{
			Comment: commentID,
			Number:  number,
		}
		if err := reportStore.CreateComment(&core.Report{ReportID: reportID}, comment); err != nil {
			c.String(500, err.Error())
			return
		}
		c.String(200, "ok")
	}
}

func hasPermission(
	c *gin.Context,
	store core.RepoStore,
	service core.SCMService,
	reportID string,
) bool {
	repo, err := store.Find(&core.Repo{ReportID: reportID})
	if err != nil {
		return false
	}
	if !repo.Private {
		return true
	}
	user, ok := request.UserFrom(c)
	if !ok {
		return false
	}
	client, err := service.Client(repo.SCM)
	if err != nil {
		return false
	}
	_, err = client.Repositories().Find(c.Request.Context(), user, repo.FullName())
	if err != nil {
		return false
	}
	return true
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

// getAll reports related to given reportID
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
