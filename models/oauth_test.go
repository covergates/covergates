package models

import (
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
)

func TestOAuth(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()

	store := &OAuthStore{DB: db}
	userStore := &UserStore{DB: db}

	userStore.Create(core.Gitea, &scm.User{
		Login: "oauth_user",
	}, &core.Token{})

	userStore.Create(core.Gitea, &scm.User{
		Login: "oauth_user2",
	}, &core.Token{})

	user, err := userStore.FindByLogin("oauth_user")
	user2, err := userStore.FindByLogin("oauth_user2")

	if err != nil {
		t.Fatal(err)
	}

	tokens := []*core.OAuthToken{
		{
			Code:  "code",
			Owner: user,
		},
		{
			Access:  "access",
			Refresh: "refresh",
			Owner:   user,
		},
	}

	for _, token := range tokens {
		if err := store.Create(token); err != nil {
			t.Fatal(err)
		}
	}

	// test create
	token, err := store.Find(&core.OAuthToken{Code: "code"})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(tokens[0], token); diff != "" {
		t.Fatal(diff)
	}

	token, err = store.Find(&core.OAuthToken{Access: "access"})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(tokens[1], token); diff != "" {
		t.Fatal(diff)
	}

	// test list

	store.Create(&core.OAuthToken{
		Code:  "code2",
		Owner: user2,
	})

	foundTokens, err := store.List(user)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(foundTokens, tokens); diff != "" {
		t.Fatal(diff)
	}

	foundTokens, err = store.List(user2)
	if err != nil {
		t.Fatal(err)
	}

	if len(foundTokens) != 1 {
		t.Fatal()
	}

	// test delete
	if err := store.Delete(&core.OAuthToken{Code: "code"}); err != nil {
		t.Fatal(err)
	}
	token, err = store.Find(&core.OAuthToken{Code: "code"})
	if err == nil {
		t.Fatal("fail to delete token")
	}

	token, err = store.Find(&core.OAuthToken{Access: "access"})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(tokens[1], token); diff != "" {
		t.Log("delete to many tokens")
		t.Fatal(diff)
	}

}
