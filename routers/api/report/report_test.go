package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/mock"
	"github.com/code-devel-cover/CodeCover/routers/api/request"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func testRequest(r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder)) {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	f(w)
}

func addFormFile(w *multipart.Writer, k, name string, r io.Reader) {
	writer, err := w.CreateFormFile(k, name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(writer, r)
	if err != nil {
		log.Fatal(err)
	}
}

func TestUpload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSCMService := mock.NewMockSCMService(ctrl)
	mockClient := mock.NewMockClient(ctrl)
	mockGitService := mock.NewMockGitService(ctrl)
	mockGitRepo := mock.NewMockGitRepository(ctrl)
	mockGitCommit := mock.NewMockGitCommit(ctrl)
	mockCoverageService := mock.NewMockCoverageService(ctrl)
	mockReportStore := mock.NewMockReportStore(ctrl)
	mockRepoStore := mock.NewMockRepoStore(ctrl)
	user := &core.User{}
	coverage := &core.CoverageReport{}
	repo := &core.Repo{
		NameSpace: "org",
		Name:      "repo",
		SCM:       core.Github,
		Branch:    "bear",
	}
	report := &core.Report{
		ReportID: "1234",
		Type:     core.ReportPerl,
		Coverage: coverage,
		Commit:   "abcdef",
		Branch:   "bear",
		Files:    []string{"a"},
	}

	mockSCMService.EXPECT().Client(
		gomock.Eq(core.Github),
	).Return(mockClient, nil)

	mockClient.EXPECT().Git().Return(mockGitService)

	mockGitService.EXPECT().GitRepository(
		gomock.Any(),
		gomock.Eq(user),
		gomock.Eq(repo.FullName()),
	).Return(mockGitRepo, nil)

	mockGitRepo.EXPECT().Commit(
		gomock.Eq(report.Commit),
	).Return(mockGitCommit, nil)

	mockGitRepo.EXPECT().ListAllFiles(
		gomock.Eq(report.Commit),
	).Return(report.Files, nil)

	mockGitCommit.EXPECT().InDefaultBranch().Return(true)

	mockRepoStore.EXPECT().Find(
		gomock.Eq(&core.Repo{
			ReportID: report.ReportID,
		}),
	).Return(repo, nil)

	mockRepoStore.EXPECT().Creator(
		gomock.Eq(repo),
	).Return(user, nil)

	mockRepoStore.EXPECT().Setting(
		gomock.Eq(repo),
	).Return(&core.RepoSetting{}, nil)

	mockCoverageService.EXPECT().Report(
		gomock.Any(),
		gomock.Eq(core.ReportPerl),
		gomock.Any(),
	).Return(coverage, nil)

	mockCoverageService.EXPECT().TrimFileNames(
		gomock.Any(),
		gomock.Eq(coverage),
		gomock.Any(),
	).Return(nil)

	mockReportStore.EXPECT().Upload(
		gomock.Eq(report),
	).Return(nil)

	r := gin.Default()
	r.POST("/reports/:id/:type", HandleUpload(
		mockSCMService,
		mockCoverageService,
		mockRepoStore,
		mockReportStore,
	))

	buffer := bytes.NewBuffer([]byte{})
	w := multipart.NewWriter(buffer)
	w.WriteField("commit", "abcdef")
	file := bytes.NewBuffer([]byte("mock"))
	addFormFile(w, "file", "cover_db.zip", file)
	w.Close()

	req, _ := http.NewRequest("POST", "/reports/1234/perl", buffer)
	req.Header.Set("Content-Type", w.FormDataContentType())
	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		rst := w.Result()
		if rst.StatusCode != 200 {
			t.Fail()
		}
	})
	// test empty commit
	req, _ = http.NewRequest("POST", "/reports/1234/perl", nil)
	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		rst := w.Result()
		if rst.StatusCode != 400 {
			t.Fail()
		}
	})
}

func TestGetRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mock.NewMockRepoStore(ctrl)
	repo := &core.Repo{
		ReportID: "1234",
	}
	store.EXPECT().Find(gomock.Eq(&core.Repo{
		ReportID: "1234",
	})).Return(repo, nil)

	r := gin.Default()
	r.GET("/reports/:id/repo", HandleRepo(store))
	req, _ := http.NewRequest("GET", "/reports/1234/repo", nil)
	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		rst := w.Result()
		if rst.StatusCode != 200 {
			t.Fail()
		}
		defer rst.Body.Close()
		data, _ := ioutil.ReadAll(rst.Body)
		rtnRepo := &core.Repo{}
		json.Unmarshal(data, rtnRepo)
		if !reflect.DeepEqual(repo, rtnRepo) {
			t.Fail()
		}
	})
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := &core.Repo{
		Branch:    "master",
		Name:      "repo",
		NameSpace: "org",
		ReportID:  "1234",
		SCM:       core.Github,
	}

	report := &core.Report{
		ReportID: "1234",
		Branch:   "master",
	}

	reportStore := mock.NewMockReportStore(ctrl)
	repoStore := mock.NewMockRepoStore(ctrl)
	service := mock.NewMockSCMService(ctrl)

	repoStore.EXPECT().Find(gomock.Eq(&core.Repo{
		ReportID: repo.ReportID,
	})).AnyTimes().Return(repo, nil)
	reportStore.EXPECT().Find(gomock.Eq(&core.Report{
		ReportID: report.ReportID,
		Branch:   repo.Branch,
	})).Return(report, nil)
	r := gin.Default()
	r.GET("/reports/:id", HandleGet(reportStore, repoStore, service))

	req, _ := http.NewRequest("GET", "/reports/1234", nil)
	query := req.URL.Query()
	query.Set("latest", "1")
	req.URL.RawQuery = query.Encode()
	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		rst := w.Result()
		if rst.StatusCode != 200 {
			t.Fail()
		}
		var reports []*core.Report
		data, _ := ioutil.ReadAll(rst.Body)
		json.Unmarshal(data, &reports)
		if len(reports) < 1 || reports[0].ReportID != "1234" {
			t.Fail()
		}
	})
}

func TestGetPrivate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := &core.Repo{
		Branch:    "master",
		Name:      "repo",
		NameSpace: "org",
		ReportID:  "1234",
		SCM:       core.Github,
		Private:   true,
	}
	reportStore := mock.NewMockReportStore(ctrl)
	repoStore := mock.NewMockRepoStore(ctrl)
	service := mock.NewMockSCMService(ctrl)
	client := mock.NewMockClient(ctrl)
	repoService := mock.NewMockRepoService(ctrl)

	repoStore.EXPECT().Find(
		gomock.Eq(&core.Repo{ReportID: repo.ReportID}),
	).AnyTimes().Return(repo, nil)

	service.EXPECT().Client(
		gomock.Eq(repo.SCM),
	).Return(client, nil)

	client.EXPECT().Repositories().Return(repoService)

	repoService.EXPECT().Find(
		gomock.Any(),
		gomock.Any(),
		gomock.Eq(repo.FullName()),
	).Return(repo, fmt.Errorf(""))

	// test if no user login
	r := gin.Default()
	r.GET("/reports/:id", HandleGet(reportStore, repoStore, service))

	req, _ := http.NewRequest("GET", "/reports/1234", nil)

	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		rst := w.Result()
		if rst.StatusCode != 401 {
			t.Fail()
		}
	})

	// test if user login but without repository access right
	r = gin.Default()
	r.Use(func(c *gin.Context) {
		request.WithUser(c, &core.User{})
	})
	r.GET("/reports/:id", HandleGet(reportStore, repoStore, service))
	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		rst := w.Result()
		if rst.StatusCode != 401 {
			t.Fail()
		}
	})

}

func TestGetNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := &core.Repo{
		ReportID: "1234",
		SCM:      core.Github,
	}

	repoStore := mock.NewMockRepoStore(ctrl)
	reportStore := mock.NewMockReportStore(ctrl)
	service := mock.NewMockSCMService(ctrl)

	repoStore.EXPECT().Find(gomock.Eq(&core.Repo{
		ReportID: repo.ReportID,
	})).AnyTimes().Return(repo, nil)

	reportStore.EXPECT().Find(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	r := gin.Default()
	r.GET("/reports/:id", HandleGet(reportStore, repoStore, service))

	req, _ := http.NewRequest("GET", "/reports/1234", nil)
	query := req.URL.Query()
	query.Set("latest", "1")
	req.URL.RawQuery = query.Encode()
	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		rst := w.Result()
		if rst.StatusCode != 404 {
			t.Fail()
		}
	})
}

func TestGetTreeMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	reportStore := mock.NewMockReportStore(ctrl)
	repoStore := mock.NewMockRepoStore(ctrl)
	chartService := mock.NewMockChartService(ctrl)
	chart := mock.NewMockChart(ctrl)
	reportID := "report_id"
	repo := &core.Repo{
		ReportID: reportID,
		Branch:   "master",
	}
	old := &core.Report{
		Coverage: &core.CoverageReport{},
		ReportID: reportID,
		Commit:   "old",
	}
	new := &core.Report{
		Coverage: &core.CoverageReport{},
		ReportID: reportID,
		Commit:   "new",
	}

	repoStore.EXPECT().Find(gomock.Eq(&core.Repo{
		ReportID: reportID,
	})).Return(repo, nil)

	reportStore.EXPECT().Find(gomock.Eq(&core.Report{
		ReportID: reportID,
		Branch:   repo.Branch,
	})).Return(old, nil)

	reportStore.EXPECT().Find(gomock.Eq(&core.Report{
		ReportID: reportID,
		Commit:   new.Commit,
	})).Return(new, nil)
	chartService.EXPECT().CoverageDiffTreeMap(
		gomock.Eq(old.Coverage),
		gomock.Eq(new.Coverage),
	).Return(chart)
	chart.EXPECT().Render(gomock.Any()).Do(
		func(w io.Writer) {
			file, _ := os.Open(filepath.Join("testdata", "treemap.svg"))
			defer file.Close()
			io.Copy(w, file)
		},
	).Return(nil)

	r := gin.Default()
	r.GET("/reports/:id/:commit/treemap", HandleGetTreeMap(
		reportStore,
		repoStore,
		chartService,
	))

	req, _ := http.NewRequest("GET", fmt.Sprintf(
		"/reports/%s/%s/treemap",
		reportID,
		new.Commit,
	), nil)
	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		rst := w.Result()
		if rst.StatusCode != 200 {
			t.Fail()
		}
		file, _ := os.Open(filepath.Join("testdata", "treemap.svg"))
		defer file.Close()
		expect, _ := ioutil.ReadAll(file)
		data, _ := ioutil.ReadAll(rst.Body)
		if len(data) <= 0 || bytes.Compare(data, expect) != 0 {
			t.Fail()
		}
	})

}
