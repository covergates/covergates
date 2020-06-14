package repo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func setRouter(
	store core.RepoStore,
) *gin.Engine {
	r := gin.Default()
	r.POST("/repo", HandleCreate(store))
	return r
}

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
	}
	store := mock.NewMockRepoStore(ctrl)
	store.EXPECT().Create(gomock.Eq(repo)).Return(nil)

	data, err := json.Marshal(repo)
	if err != nil {
		t.Error(err)
		return
	}
	read := bytes.NewReader(data)
	req, _ := http.NewRequest("POST", "/repo", read)
	r := setRouter(store)
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
