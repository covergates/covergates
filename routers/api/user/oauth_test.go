package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/mock"
	"github.com/covergates/covergates/routers/api/request"
	"github.com/covergates/covergates/routers/api/user"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func testRequest(r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder)) {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	f(w)
}

func TestTokenCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := &core.User{
		Login: "user",
	}
	mockToken := &core.OAuthToken{
		Name:   "token",
		Access: "access",
	}

	mockService := mock.NewMockOAuthService(ctrl)
	mockService.EXPECT().CreateToken(
		gomock.Any(),
		gomock.Eq(mockToken.Name),
	).Return(
		mockToken, nil,
	)

	mockService.EXPECT().WithUser(
		gomock.Any(),
		gomock.Eq(mockUser),
	).Return(context.Background())

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		request.WithUser(c, mockUser)
	})

	r.POST("/tokens", user.HandleCreateToken(mockService))

	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	w.WriteField("name", mockToken.Name)
	w.Close()
	request, _ := http.NewRequest("POST", "/tokens", buf)
	request.Header.Set("Content-Type", w.FormDataContentType())
	testRequest(r, request, func(w *httptest.ResponseRecorder) {
		respond := w.Result()
		if respond.StatusCode != 200 {
			t.Fatal()
		}
		defer respond.Body.Close()
		data, _ := ioutil.ReadAll(respond.Body)
		if diff := cmp.Diff(string(data), mockToken.Access); diff != "" {
			t.Fatal(diff)
		}
	})
}

func TestListTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockOAuthService(ctrl)

	c1 := mockService.EXPECT().ListTokens(
		gomock.Any(),
	).Return(
		[]*core.OAuthToken{
			{
				ID:   1,
				Name: "token1",
			},
		},
		nil,
	)

	mockService.EXPECT().ListTokens(
		gomock.Any(),
	).Return(nil, errors.New("")).After(c1)

	mockService.EXPECT().WithUser(
		gomock.Any(), gomock.Any(),
	).Return(context.Background()).AnyTimes()

	t.Run("should check user", func(t *testing.T) {
		r := gin.Default()
		r.GET("/tokens", user.HandleListTokens(mockService))
		request, _ := http.NewRequest("GET", "/tokens", nil)
		testRequest(r, request, func(w *httptest.ResponseRecorder) {
			respond := w.Result()
			if respond.StatusCode != 401 {
				t.Fatal()
			}
		})
	})

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		request.WithUser(c, &core.User{Login: "login"})
	})
	r.GET("/tokens", user.HandleListTokens(mockService))
	t.Run("should get tokens", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/tokens", nil)
		testRequest(r, request, func(w *httptest.ResponseRecorder) {
			respond := w.Result()
			if respond.StatusCode != 200 {
				t.Fatal()
			}
			defer respond.Body.Close()
			data, _ := ioutil.ReadAll(respond.Body)
			var tokens []*user.Token
			json.Unmarshal(data, &tokens)
			if len(tokens) < 1 || tokens[0].Name != "token1" {
				t.Fatal()
			}
		})
	})
	t.Run("should return 500 when error", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/tokens", nil)
		testRequest(r, request, func(w *httptest.ResponseRecorder) {
			respond := w.Result()
			if respond.StatusCode != 500 {
				t.Fatal()
			}
		})
	})
}

func TestDeleteToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockOAuthService(ctrl)
	mockStore := mock.NewMockOAuthStore(ctrl)

	t.Run("should check user", func(t *testing.T) {
		r := gin.Default()
		r.DELETE("/tokens/:id", user.HandleDeleteToken(mockService, mockStore))
		request, _ := http.NewRequest("DELETE", "/tokens/1", nil)
		testRequest(r, request, func(w *httptest.ResponseRecorder) {
			response := w.Result()
			if response.StatusCode != 401 {
				t.Fatal()
			}
		})
	})

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		request.WithUser(c, &core.User{Login: "login"})
	})
	r.DELETE("/tokens/:id", user.HandleDeleteToken(mockService, mockStore))

	t.Run("should check id", func(t *testing.T) {
		request, _ := http.NewRequest("DELETE", "/tokens/bear", nil)
		testRequest(r, request, func(w *httptest.ResponseRecorder) {
			response := w.Result()
			if response.StatusCode != 400 {
				t.Fatal()
			}
		})
	})

	t.Run("should return 500 when not found", func(t *testing.T) {
		mockStore.EXPECT().Find(gomock.Any()).Return(nil, errors.New(""))
		request, _ := http.NewRequest("DELETE", "/tokens/1", nil)
		testRequest(r, request, func(w *httptest.ResponseRecorder) {
			response := w.Result()
			if response.StatusCode != 500 {
				t.Fatal()
			}
		})
	})

	mockService.EXPECT().WithUser(
		gomock.Any(), gomock.Any(),
	).Return(context.Background()).AnyTimes()

	t.Run("should return 500 when fail deleting", func(t *testing.T) {
		token := &core.OAuthToken{}
		mockStore.EXPECT().Find(gomock.Any()).Return(token, nil)
		mockService.EXPECT().DeleteToken(gomock.Any(), token).Return(errors.New(""))
		request, _ := http.NewRequest("DELETE", "/tokens/1", nil)
		testRequest(r, request, func(w *httptest.ResponseRecorder) {
			resposne := w.Result()
			if resposne.StatusCode != 500 {
				t.Fatal()
			}
		})
	})

	t.Run("should return deleted token", func(t *testing.T) {
		token := &core.OAuthToken{ID: 1, Name: "token"}
		mockStore.EXPECT().Find(
			&core.OAuthToken{ID: token.ID},
		).Return(token, nil)
		mockService.EXPECT().DeleteToken(gomock.Any(), token).Return(nil)
		request, _ := http.NewRequest("DELETE", "/tokens/1", nil)
		testRequest(r, request, func(w *httptest.ResponseRecorder) {
			response := w.Result()
			if response.StatusCode != 200 {
				t.Fatal()
			}
			var token user.Token
			defer response.Body.Close()
			data, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(data, &token); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(user.Token{ID: 1, Name: "token"}, token); diff != "" {
				t.Fatal(diff)
			}
		})
	})

}
