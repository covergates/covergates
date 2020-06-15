package scm

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/mock"
	"github.com/code-devel-cover/CodeCover/routers/api/request"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func setRouter(
	repoService core.RepoService,
) *gin.Engine {
	r := gin.Default()
	return r
}

func testRequest(r *gin.Engine, req *http.Request, f func(*httptest.ResponseRecorder)) {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	f(w)
}

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	user := &core.User{Login: "user"}

	scmService := mock.NewMockSCMService(ctrl)
	client := mock.NewMockClient(ctrl)
	repoService := mock.NewMockRepoService(ctrl)
	scmService.EXPECT().Client(gomock.Eq(core.Github)).Return(client, nil)
	client.EXPECT().Repositories().Return(repoService)
	repoService.EXPECT().List(gomock.Any(), gomock.Eq(user)).Return(
		[]*core.Repo{
			{
				Name: "repo1",
			},
		},
		nil,
	)

	r := setRouter(repoService)
	r.Use(func(c *gin.Context) {
		request.WithUser(c, user)
	})
	r.GET("/scm/:scm/repos", HandleListSCM(scmService))
	req, _ := http.NewRequest("GET", "/scm/github/repos", nil)
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
		var repositories []*core.Repo
		if err := json.Unmarshal(data, &repositories); err != nil {
			t.Error(err)
			return
		}
		if len(repositories) < 1 || repositories[0].Name != "repo1" {
			t.Fail()
		}
	})
	// test no login
	r = setRouter(repoService)
	r.GET("/scm/:scm/repos", HandleListSCM(scmService))
	req, _ = http.NewRequest("GET", "/scm/github/repos", nil)
	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		rst := w.Result()
		if rst.StatusCode != 500 {
			t.Fail()
		}
	})
}
