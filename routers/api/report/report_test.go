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

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/mock"
	"github.com/covergates/covergates/modules/charts"
	"github.com/covergates/covergates/routers/api/request"
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

func encodeJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
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

func createForm(writer io.Writer, m map[string]string) *multipart.Writer {
	w := multipart.NewWriter(writer)
	for k, v := range m {
		w.WriteField(k, v)
	}
	return w
}

func TestUpload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCoverageService := mock.NewMockCoverageService(ctrl)
	mockReportStore := mock.NewMockReportStore(ctrl)

	t.Run("basic", func(t *testing.T) {
		coverage := &core.CoverageReport{
			Type: core.ReportPerl,
		}
		report := &core.Report{
			ReportID: "1234",
			Coverages: []*core.CoverageReport{
				coverage,
			},
			Commit:    "abcdef",
			Reference: "bear",
			Files:     []string{"a"},
		}

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
		r.Use(func(c *gin.Context) {
			WithSetting(c, &core.RepoSetting{})
		})
		r.POST("/reports/:id", HandleUpload(
			mockCoverageService,
			mockReportStore,
		))
		buffer := bytes.NewBuffer([]byte{})
		w := createForm(
			buffer,
			map[string]string{
				"commit": "abcdef",
				"type":   "perl",
				"ref":    "bear",
				"files":  encodeJSON(report.Files),
			},
		)
		addFormFile(w, "file", "cover_db.zip", bytes.NewBuffer([]byte("mock")))
		w.Close()

		req, _ := http.NewRequest("POST", "/reports/1234", buffer)
		req.Header.Set("Content-Type", w.FormDataContentType())
		testRequest(r, req, func(w *httptest.ResponseRecorder) {
			rst := w.Result()
			if rst.StatusCode != 200 {
				t.Fail()
			}
		})
	})

	t.Run("test empty post", func(t *testing.T) {
		r := gin.Default()
		r.Use(func(c *gin.Context) {
			WithSetting(c, &core.RepoSetting{})
		})
		r.POST("/reports/:id", HandleUpload(
			mockCoverageService,
			mockReportStore,
		))
		req, _ := http.NewRequest("POST", "/reports/1234", nil)
		testRequest(r, req, func(w *httptest.ResponseRecorder) {
			rst := w.Result()
			if rst.StatusCode != 400 {
				t.Fail()
			}
		})
	})
}

func TestProtectReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoStore := mock.NewMockRepoStore(ctrl)

	t.Run("test protected report", func(t *testing.T) {
		r := gin.Default()
		r.POST("/",
			func(c *gin.Context) {
				WithSetting(c, &core.RepoSetting{
					Protected: true,
				})
			},
			ProtectReport(func(c *gin.Context) {
				c.String(401, "")
				c.Abort()
			}, mockRepoStore),
		)
		request, _ := http.NewRequest("POST", "/", nil)
		testRequest(r, request, func(w *httptest.ResponseRecorder) {
			response := w.Result()
			if response.StatusCode != 401 {
				t.Fatal()
			}
		})
	})

	t.Run("test unprotected report", func(t *testing.T) {
		r := gin.Default()
		r.POST("/",
			func(c *gin.Context) {
				WithSetting(c, &core.RepoSetting{
					Protected: false,
				})
			},
			ProtectReport(func(c *gin.Context) {
				c.String(401, "")
				c.Abort()
			}, mockRepoStore),
		)
		request, _ := http.NewRequest("POST", "/", nil)
		testRequest(r, request, func(w *httptest.ResponseRecorder) {
			response := w.Result()
			if response.StatusCode != 200 {
				t.Fatal()
			}
		})
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
		ReportID:  "1234",
		Reference: "master",
	}

	reportStore := mock.NewMockReportStore(ctrl)
	repoStore := mock.NewMockRepoStore(ctrl)
	service := mock.NewMockSCMService(ctrl)

	repoStore.EXPECT().Find(gomock.Eq(&core.Repo{
		ReportID: repo.ReportID,
	})).AnyTimes().Return(repo, nil)
	reportStore.EXPECT().Find(gomock.Eq(&core.Report{
		ReportID:  report.ReportID,
		Reference: repo.Branch,
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
	repoService := mock.NewMockGitRepoService(ctrl)

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
		Coverages: []*core.CoverageReport{
			{
				Type: core.ReportGo,
			},
		},
		ReportID:  reportID,
		Reference: "master",
		Commit:    "old",
	}
	new := &core.Report{
		Coverages: []*core.CoverageReport{
			{
				Type: core.ReportGo,
			},
		},
		ReportID:  reportID,
		Reference: "new",
		Commit:    "new",
	}

	repoStore.EXPECT().Find(gomock.Eq(&core.Repo{
		ReportID: reportID,
	})).Return(repo, nil)

	reportStore.EXPECT().Find(gomock.Eq(
		&core.Report{
			Commit:   new.Reference,
			ReportID: new.ReportID,
		},
	)).Return(nil, fmt.Errorf(""))

	reportStore.EXPECT().Find(gomock.Eq(&core.Report{
		ReportID:  reportID,
		Reference: repo.Branch,
	})).Return(old, nil)

	reportStore.EXPECT().Find(gomock.Eq(&core.Report{
		ReportID:  reportID,
		Reference: new.Reference,
	})).Return(new, nil)
	chartService.EXPECT().CoverageDiffTreeMap(
		gomock.Eq(old),
		gomock.Eq(new),
	).Return(chart)
	chart.EXPECT().Render(gomock.Any()).Do(
		func(w io.Writer) {
			file, _ := os.Open(filepath.Join("testdata", "treemap.svg"))
			defer file.Close()
			io.Copy(w, file)
		},
	).Return(nil)

	r := gin.Default()
	r.GET("/reports/:id/treemap/*ref", HandleGetTreeMap(
		reportStore,
		repoStore,
		chartService,
	))

	req, _ := http.NewRequest("GET", fmt.Sprintf(
		"/reports/%s/treemap/%s",
		reportID,
		new.Reference,
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

func TestGetCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepoStore(ctrl)
	mockReport := mock.NewMockReportStore(ctrl)
	mockChart := mock.NewMockChartService(ctrl)

	repo := &core.Repo{
		Branch:   "master",
		ReportID: "report_id",
	}
	report := &core.Report{
		ReportID: "report_id",
	}

	mockRepo.EXPECT().Find(
		gomock.Eq(&core.Repo{ReportID: repo.ReportID}),
	).Return(repo, nil)
	mockReport.EXPECT().Find(
		gomock.Eq(&core.Report{ReportID: repo.ReportID, Reference: repo.Branch}),
	).Return(report, nil)
	mockChart.EXPECT().RepoCard(
		gomock.Eq(repo),
		gomock.Eq(report),
	).Return(charts.NewRepoCard(repo, report))

	r := gin.Default()
	r.GET("/reports/:id/card", HandleGetCard(
		mockRepo,
		mockReport,
		mockChart,
	))

	req, _ := http.NewRequest("GET", fmt.Sprintf(
		"/reports/%s/card",
		repo.ReportID,
	), nil)

	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		result := w.Result()
		if result.StatusCode != 200 {
			t.Fatal()
		}
	})
}
