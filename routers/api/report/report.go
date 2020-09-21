package report

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

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
// @Param ref formData string false "ref"
// @Param files formData string false "files list of the repository"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "error message"
// @Router /reports/{id} [post]
func HandleUpload(
	coverageService core.CoverageService,
	reportStore core.ReportStore,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := c.GetPostForm("type"); !ok {
			c.String(400, "must have report type")
			return
		}

		if _, ok := c.GetPostForm("commit"); !ok {
			c.String(400, "must have commit SHA")
			return
		}

		reportID := c.Param("id")
		ref := c.PostForm("ref")
		reportType := core.ReportType(c.PostForm("type"))
		commit := c.PostForm("commit")
		ctx := c.Request.Context()

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
		coverage, err := loadCoverageReport(
			ctx,
			coverageService,
			reportType,
			reader,
			MustGetSetting(c),
		)
		if err != nil {
			log.Error(err)
			c.String(500, err.Error())
			return
		}

		report := &core.Report{
			ReportID: reportID,
			Coverages: []*core.CoverageReport{
				coverage,
			},
			Files:     files,
			Reference: ref,
			Commit:    commit,
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
	Latest bool   `form:"latest"`
	Ref    string `form:"ref"`
}

// HandleGet for the report id
// @Summary get reports for the report id
// @Tags Report
// @Param id path string true "report id"
// @Param latest query bool false "get only the latest report"
// @Param ref query string false "get report for git ref"
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
		var err error
		var reports []*core.Report
		if option.Latest && option.Ref == "" {
			var report *core.Report
			if report, err = getLatest(reportStore, repoStore, reportID); err == nil {
				reports = []*core.Report{report}
			}
		} else if option.Latest && option.Ref != "" {
			var report *core.Report
			if report, err = getRef(reportStore, reportID, option.Ref); err == nil {
				reports = []*core.Report{report}
			}
		} else if option.Ref != "" {
			reports, err = reportStore.List(reportID, option.Ref)
		} else {
			reports, err = getAll(reportStore, reportID)
		}
		if err != nil {
			c.Error(err)
			c.JSON(404, []*core.Report{})
			return
		}
		c.JSON(200, reports)
	}
}

// HandleGetTreeMap for coverage difference with main branch
// @Summary Get coverage difference treemap with main branch
// @Tags Report
// @Produce image/svg+xml
// @Param id path string true "report id"
// @param source path string true "source branch"
// @Success 200 {object} string "treemap svg"
// @Router /reports/{id}/treemap/{ref} [get]
func HandleGetTreeMap(
	reportStore core.ReportStore,
	repoStore core.RepoStore,
	chartService core.ChartService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		reportID := c.Param("id")
		ref := strings.Trim(c.Param("ref"), "/")
		new, err := getRef(reportStore, reportID, ref)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		old, err := getLatest(reportStore, repoStore, reportID)
		if err != nil {
			old = &core.Report{
				Coverages: []*core.CoverageReport{},
			}
		}
		chart := chartService.CoverageDiffTreeMap(old, new)
		buffer := bytes.NewBuffer([]byte{})
		if err := chart.Render(buffer); err != nil {
			c.String(500, err.Error())
			return
		}
		c.Header("Cache-Control", "max-age=600")
		c.Data(200, "image/svg+xml", buffer.Bytes())
		return
	}
}

// HandleGetCard of the repository status
// @Summary Get status card of the repository
// @Tags Report
// @Produce image/svg+xml
// @Param id path string true "report id"
// @Success 200 {object} string "treemap svg"
// @Router /reports/{id}/card [get]
func HandleGetCard(
	repoStore core.RepoStore,
	reportStore core.ReportStore,
	chartService core.ChartService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		reportID := c.Param("id")
		repo, err := repoStore.Find(&core.Repo{ReportID: reportID})
		if err != nil {
			c.String(404, "repository not found")
			return
		}
		report, err := reportStore.Find(&core.Report{ReportID: reportID, Reference: repo.Branch})
		if err != nil {
			c.String(404, "report not found")
			return
		}
		chart := chartService.RepoCard(repo, report)
		buffer := bytes.NewBuffer([]byte{})
		if err := chart.Render(buffer); err != nil {
			c.String(500, err.Error())
			return
		}
		c.Header("Cache-Control", "max-age=600")
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
			if source, err = reportStore.Find(&core.Report{ReportID: reportID, Reference: pr.Source}); err != nil {
				c.String(500, err.Error())
				return
			}
		}
		target, err := reportStore.Find(&core.Report{ReportID: reportID, Reference: pr.Target})
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
			source.Reference,
			target.Reference,
		))

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
		ReportID:  reportID,
		Reference: repo.Branch,
	})
}

func getRef(store core.ReportStore, reportID, ref string) (*core.Report, error) {
	var report *core.Report
	var err error
	seed := &core.Report{ReportID: reportID, Commit: ref}
	if report, err = store.Find(seed); err == nil {
		return report, err
	}
	seed = &core.Report{ReportID: reportID, Reference: ref}
	if report, err = store.Find(seed); err == nil {
		return report, err
	}
	return nil, err
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
	coverage.Type = reportType
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
