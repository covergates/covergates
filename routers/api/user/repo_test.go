package user_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/mock"
	"github.com/covergates/covergates/routers/api/request"
	"github.com/covergates/covergates/routers/api/user"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestHandleSynchronizeRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// data
	u := &core.User{
		Login: "test",
	}

	// mock
	service := mock.NewMockRepoService(ctrl)
	service.EXPECT().Synchronize(gomock.Any(), gomock.Eq(u)).Times(1).Return(nil)
	service.EXPECT().Synchronize(gomock.Any(), gomock.Eq(u)).Times(1).Return(errors.New(""))

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		request.WithUser(c, u)
	})
	r.PATCH("/user/repos", user.HandleSynchronizeRepo(service))

	req, _ := http.NewRequest("PATCH", "/user/repos", nil)

	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		result := w.Result()
		if result.StatusCode != 200 {
			t.Fatal()
		}
	})

	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		result := w.Result()
		if result.StatusCode != 500 {
			t.Fatal()
		}
	})
}

func TestHandleListRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// data
	u := &core.User{}

	// mock
	store := mock.NewMockUserStore(ctrl)
	store.EXPECT().ListRepositories(gomock.Eq(u)).Return([]*core.Repo{}, nil)
	store.EXPECT().ListRepositories(gomock.Eq(u)).Return(nil, errors.New(""))

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		request.WithUser(c, u)
	})
	r.GET("/user/repos", user.HandleListRepo(store))

	req, _ := http.NewRequest("GET", "/user/repos", nil)

	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		result := w.Result()
		if result.StatusCode != 200 {
			t.Fatal()
		}
	})

	testRequest(r, req, func(w *httptest.ResponseRecorder) {
		result := w.Result()
		if result.StatusCode != 500 {
			t.Fatal()
		}
		defer result.Body.Close()
		data, _ := ioutil.ReadAll(result.Body)
		var repos []*core.Repo
		if err := json.Unmarshal(data, &repos); err != nil {
			t.Fatal(err)
		}
		if len(repos) != 0 {
			t.Fatal()
		}
	})
}
