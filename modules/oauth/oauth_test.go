package oauth_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/models"
	"github.com/covergates/covergates/modules/oauth"
	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var service *oauth.Service
var userStore core.UserStore
var conf *config.Config

func mockUsers(store core.UserStore) {
	store.Create(core.Gitea, &scm.User{Login: "user1"}, &core.Token{})
	store.Create(core.Gitea, &scm.User{Login: "user2"}, &core.Token{})
}

func TestMain(m *testing.M) {
	log.SetReportCaller(true)
	cwd, _ := os.Getwd()
	tempFile, err := ioutil.TempFile(cwd, "*.db")
	if err != nil {
		log.Fatal(err)
	}
	tempFile.Close()
	x, err := gorm.Open(sqlite.Open(tempFile.Name()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	dbService := models.NewDatabaseService(x)
	userStore = &models.UserStore{DB: dbService}
	oauthStore := &models.OAuthStore{DB: dbService}
	dbService.Migrate()
	mockUsers(userStore)

	conf = &config.Config{}
	service = oauth.NewService(conf, oauthStore, userStore)
	exit := m.Run()
	os.Remove(tempFile.Name())
	os.Exit(exit)
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	if _, err := service.CreateToken(ctx, ""); err == nil || err != oauth.ErrTokenOwnerNotFound {
		t.Fatal("should check token owner in contex")
	}

	user, err := userStore.FindByLogin("user1")
	if err != nil {
		t.Fatal(err)
	}

	ctx = service.WithUser(ctx, user)

	token, err := service.CreateToken(ctx, "test_token")
	if err != nil {
		t.Fatal(err)
	}

	if token.Access == "" || token.Name != "test_token" {
		t.Fatal()
	}

	if diff := cmp.Diff(user, token.Owner); diff != "" {
		t.Fatal(diff)
	}
}

func TestValidate(t *testing.T) {

	user, err := userStore.FindByLogin("user1")
	if err != nil {
		t.Fatal(err)
	}

	ctx := service.WithUser(context.Background(), user)

	token, err := service.CreateToken(ctx, "validate_token")

	if err != nil {
		t.Fatal(err)
	}

	request, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("%s?access_token=%s", conf.Server.Addr, token.Access),
		nil,
	)

	tokenOwner, err := service.Validate(request)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(user, tokenOwner); diff != "" {
		t.Fatal(diff)
	}

	request, _ = http.NewRequest(
		"GET",
		fmt.Sprintf("%s?access_token=123", conf.Server.Addr),
		nil,
	)
	if _, err := service.Validate(request); err == nil {
		t.Fatal("should return err for invalid token")
	}
}

func TestDelete(t *testing.T) {
	user, err := userStore.FindByLogin("user1")
	if err != nil {
		t.Fatal(err)
	}

	ctx := service.WithUser(context.Background(), user)

	token, err := service.CreateToken(ctx, "delete_token")

	if err != nil {
		t.Fatal(err)
	}

	if err := service.DeleteToken(ctx, token); err != nil {
		t.Fatal(err)
	}

	request, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("%s?access_token=%s", conf.Server.Addr, token.Access),
		nil,
	)
	if _, err := service.Validate(request); err == nil {
		t.Fatal("should return err for deleted token")
	}

	token, err = service.CreateToken(ctx, "user1_token")
	if err != nil {
		t.Fatal(err)
	}

	user2, err := userStore.FindByLogin("user2")
	if err != nil {
		t.Fatal(err)
	}
	ctx = service.WithUser(context.Background(), user2)
	if err := service.DeleteToken(ctx, token); err == nil {
		t.Fatal("user2 cannot delete user1's token")
	}
}

func TestList(t *testing.T) {
	user, err := userStore.FindByLogin("user2")
	if err != nil {
		t.Fatal(err)
	}

	names := []string{"token1", "token2"}

	ctx := service.WithUser(context.Background(), user)
	for _, name := range names {
		service.CreateToken(ctx, name)
	}

	tokens, err := service.ListTokens(ctx)
	if err != nil {
		t.Fatal(err)
	}

	tokenNames := make([]string, len(tokens))
	for i, token := range tokens {
		tokenNames[i] = token.Name
	}
	if diff := cmp.Diff(names, tokenNames); diff != "" {
		t.Fatal(diff)
	}
}
