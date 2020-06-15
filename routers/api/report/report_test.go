package report

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
)

func setupRouter(
	coverageService core.CoverageService,
	reportStore core.ReportStore,
) *gin.Engine {
	r := gin.Default()
	r.POST("/report/:id/:type", HandleUpload(coverageService, reportStore))
	return r
}

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
	mockCoverageService := mock.NewMockCoverageService(ctrl)
	mockReportStore := mock.NewMockReportStore(ctrl)
	coverage := &core.CoverageReport{}
	mockCoverageService.EXPECT().Report(
		gomock.Any(),
		gomock.Eq(core.ReportPerl),
		gomock.Any(),
	).Return(coverage, nil)
	mockReportStore.EXPECT().Upload(gomock.Eq(
		&core.Report{
			ReportID: "1234",
			Type:     core.ReportPerl,
			Coverage: coverage,
			Commit:   "abcdef",
		},
	)).Return(nil)
	r := setupRouter(mockCoverageService, mockReportStore)

	buffer := bytes.NewBuffer([]byte{})
	w := multipart.NewWriter(buffer)
	w.WriteField("commit", "abcdef")
	file := bytes.NewBuffer([]byte("mock"))
	addFormFile(w, "file", "cover_db.zip", file)
	w.Close()

	req, _ := http.NewRequest("POST", "/report/1234/perl", buffer)
	req.Header.Set("Content-Type", w.FormDataContentType())
	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		rst := w.Result()
		if rst.StatusCode != 200 {
			t.Fail()
		}
	})
	// test empty commit
	req, _ = http.NewRequest("POST", "/report/1234/perl", nil)
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
	r.GET("/report/:id/repo", HandleRepo(store))
	req, _ := http.NewRequest("GET", "/report/1234/repo", nil)
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
