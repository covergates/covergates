package models

import (
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func TestUserCreate(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := &UserStore{
		DB: db,
	}
	user1 := &scm.User{
		Login:  "create1",
		Name:   "create1",
		Email:  "create1@gmail.com",
		Avatar: "http://avatar",
	}
	token1 := &core.Token{}
	if err := store.Create(core.Github, user1, token1); err != nil {
		t.Error(err)
		return
	}
	if err := store.Create(core.Gitea, user1, token1); err == nil {
		t.Fail()
	}
}

func TestUserFind(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := &UserStore{
		DB: db,
	}
	session := store.DB.Session()
	user1 := &User{
		GithubLogin: "user1",
		GithubEmail: "user1@gmail.com",
		Email:       "user1@gmail.com",
		Login:       "user1",
	}
	if err := session.Create(user1).Error; err != nil {
		t.Error(err)
		return
	}
	scmUser1 := scm.User{
		Login: "user1",
		Email: "user1@gmail.com",
	}
	coreUser1, err := store.Find(core.Github, &scmUser1)
	if err != nil {
		t.Error(err)
		return
	}
	if coreUser1.Login != "user1" || coreUser1.GithubEmail != "user1@gmail.com" {
		t.Fail()
	}
	scmUser2 := scm.User{
		Login: "user2",
		Email: "user2@gmail.com",
	}
	coreUser2, err := store.Find(core.Github, &scmUser2)
	if err == nil || !gorm.IsRecordNotFoundError(err) {
		t.Fail()
	}
	if coreUser2 != nil {
		t.Fail()
	}
}

func TestUserBind(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := &UserStore{
		DB: db,
	}
	user := &User{
		GithubLogin: "bindgithub",
		GithubEmail: "bindgithub@gmail.com",
		Email:       "bindgithub@gmail.com",
		Login:       "bindgithub",
	}

	giteaUser := &scm.User{
		Email: "bindgitea@gmail.com",
		Login: "bindgitea",
	}

	expectUser := &core.User{
		Login:       user.Login,
		Email:       user.GithubEmail,
		GithubLogin: user.GithubLogin,
		GithubEmail: user.GithubEmail,
		GiteaLogin:  giteaUser.Login,
		GiteaEmail:  giteaUser.Email,
	}

	_, err := store.Bind(core.Gitea, user.toCoreUser(), giteaUser, &core.Token{})
	if err == nil {
		log.Info("could not bind with inexistent user")
		t.Fail()
		return
	}

	err = store.Create(core.Github, &scm.User{
		Login: user.GithubLogin,
		Email: user.GithubEmail,
	}, &core.Token{})

	if err != nil {
		t.Fatal(err)
	}

	newUser, err := store.Bind(core.Gitea, user.toCoreUser(), giteaUser, &core.Token{})
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}

	if diff := cmp.Diff(newUser, expectUser); diff != "" {
		t.Log(diff)
		t.Fail()
	}
}
