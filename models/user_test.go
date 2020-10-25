package models

import (
	"errors"
	"testing"
	"time"

	"github.com/covergates/covergates/core"
	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
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

	trans := cmp.Transformer("", func(in *core.User) *core.User {
		in.GitLabExpire = time.Unix(0, 0)
		in.BitbucketExpire = time.Unix(0, 0)
		return in
	})

	if diff := cmp.Diff(newUser, expectUser, trans); diff != "" {
		t.Log(diff)
		t.Fail()
	}
}

func TestUserBindExist(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := &UserStore{
		DB: db,
	}

	githubUser := &scm.User{
		Login: "exists_github",
		Email: "exists_github@mail",
	}

	giteaUser := &scm.User{
		Login: "exists_gitea",
		Email: "exists_gitea@mail",
	}

	err := store.Create(core.Github, githubUser, &core.Token{})
	if err != nil {
		t.Fatal(err)
	}
	err = store.Create(core.Gitea, giteaUser, &core.Token{})
	if err != nil {
		t.Fatal(err)
	}

	user1, err := store.Find(core.Github, githubUser)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Bind(core.Gitea, user1, giteaUser, &core.Token{}); err != errUserExist {
		t.Fail()
	}
}

func newUser(t *testing.T, store *UserStore, provider core.SCMProvider, login string) *core.User {
	if err := store.Create(
		provider,
		&scm.User{Login: login},
		&core.Token{},
	); err != nil {
		t.Fatal(err)
	}

	user, err := store.FindByLogin(login)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func newRepos(t *testing.T, db core.DatabaseService, repos []*Repo) []*core.Repo {
	coreRepos := make([]*core.Repo, len(repos))
	for i, repo := range repos {
		if err := db.Session().Create(repo).Error; err != nil {
			t.Fatal(err)
		}
		coreRepos[i] = repo.ToCoreRepo()
	}
	return coreRepos
}

func TestUserRepositories(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()

	userStore := &UserStore{DB: db}

	repos := newRepos(
		t, db,
		[]*Repo{
			{
				URL: "http://testuser/repo1",
				SCM: string(core.Github),
			},
			{
				URL: "http://testuser/repo2",
				SCM: string(core.Github),
			},
		},
	)

	user := newUser(t, userStore, core.Gitea, "test_repo_user1")

	// test new create
	if err := userStore.UpdateRepositories(user, repos); err != nil {
		t.Fatal(err)
	}

	result, err := userStore.ListRepositories(user)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(repos, result); diff != "" {
		t.Fatal(diff)
	}

	// test update, remove repository
	if err := userStore.UpdateRepositories(user, repos[0:1]); err != nil {
		t.Fatal()
	}

	result, err = userStore.ListRepositories(user)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(repos[0:1], result); diff != "" {
		t.Fatal(diff)
	}

	// test add inexistent repository
	if err := userStore.UpdateRepositories(user, append(repos, &core.Repo{URL: "fake"})); err == nil {
		t.Fatal()
	}

	// test second user
	user2 := newUser(t, userStore, core.Github, "test_repo_user2")
	repos2 := newRepos(t, db, []*Repo{
		{
			URL: "http://testuser/github/repo1",
			SCM: string(core.Github),
		},
		{
			URL: "http://testuser/github/repo2",
			SCM: string(core.Github),
		},
	})
	if err := userStore.UpdateRepositories(user2, repos2); err != nil {
		t.Fatal(err)
	}
	result, err = userStore.ListRepositories(user2)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(repos2, result); diff != "" {
		t.Fatal(diff)
	}
}
