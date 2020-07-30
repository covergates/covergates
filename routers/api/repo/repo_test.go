package repo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/mock"
	"github.com/covergates/covergates/routers/api/request"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func testRequest(r *gin.Engine, req *http.Request, f func(*httptest.ResponseRecorder)) {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	f(w)
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := &core.Repo{
		URL:       "http://gitea/org/repo",
		NameSpace: "org",
		Name:      "repo",
		SCM:       core.Gitea,
		Branch:    "master",
	}
	user := &core.User{}
	store := mock.NewMockRepoStore(ctrl)
	store.EXPECT().Create(gomock.Eq(repo), gomock.Eq(user)).Return(nil)
	service := mock.NewMockSCMService(ctrl)

	data, err := json.Marshal(repo)
	if err != nil {
		t.Error(err)
		return
	}
	read := bytes.NewReader(data)
	req, _ := http.NewRequest("POST", "/repo", read)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		request.WithUser(c, user)
	})
	r.POST("/repo", HandleCreate(store, service))
	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		rst := w.Result()
		if rst.StatusCode != 200 {
			t.Fail()
			return
		}
		data, err := ioutil.ReadAll(rst.Body)
		if err != nil {
			t.Error(err)
			return
		}
		rstRepo := &core.Repo{}
		json.Unmarshal(data, rstRepo)
		if !reflect.DeepEqual(repo, rstRepo) {
			t.Fail()
		}
	})
}

func TestListSCM(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &core.User{}
	scmRepos := []*core.Repo{
		{
			Name: "repo1",
			URL:  "url1",
		},
		{
			Name: "repo2",
			URL:  "url2",
		},
	}
	storeRepos := []*core.Repo{
		{
			URL:      "url2",
			ReportID: "report_id",
		},
	}
	urls := make([]string, len(scmRepos))
	for i, repo := range scmRepos {
		urls[i] = repo.URL
	}

	mockService := mock.NewMockSCMService(ctrl)
	mockClient := mock.NewMockClient(ctrl)
	mockRepoService := mock.NewMockRepoService(ctrl)
	mockStore := mock.NewMockRepoStore(ctrl)

	mockService.EXPECT().Client(gomock.Eq(core.Github)).Return(mockClient, nil)
	mockClient.EXPECT().Repositories().Return(mockRepoService)
	mockRepoService.EXPECT().List(gomock.Any(), gomock.Eq(user)).Return(scmRepos, nil)
	mockStore.EXPECT().Finds(gomock.Eq(urls)).Return(storeRepos, nil)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		request.WithUser(c, user)
	})
	r.GET("/repos/:scm", HandleListSCM(mockService, mockStore))

	req, _ := http.NewRequest("GET", "/repos/github", nil)
	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		rst := w.Result()
		if rst.StatusCode != 200 {
			t.Fail()
			return
		}
		data, _ := ioutil.ReadAll(rst.Body)
		var repos []*core.Repo
		json.Unmarshal(data, &repos)
		if len(repos) < 2 {
			t.Fail()
			return
		}
		if repos[0].ReportID != "" {
			t.Fail()
		}
		if repos[1].ReportID != "report_id" {
			t.Fail()
		}
	})
}
