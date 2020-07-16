package models

import (
	"reflect"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
)

func TestRepoFind(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	session := db.Session()
	repo := &Repo{
		Name:      "test_repo_find",
		NameSpace: "gitea",
		SCM:       string(core.Gitea),
		URL:       "http://gitea/test_repo_find",
	}
	session.Create(repo)
	if repo.ID <= 0 {
		t.Fail()
	}
	store := &RepoStore{DB: db}
	rst, err := store.Find(&core.Repo{
		Name:      "test_repo_find",
		NameSpace: "gitea",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if rst.URL != repo.URL {
		t.Fail()
	}
	rst, err = store.Find(&core.Repo{ID: repo.ID})
	if err != nil {
		t.Error(err)
		return
	}
	if rst.Name != "test_repo_find" {
		t.Fail()
	}
}

func TestRepoFinds(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	session := db.Session()
	names := []string{"a", "b", "c"}
	urls := make([]string, len(names))
	for i, name := range names {
		urls[i] = "http://gitea/finds" + name
		session.Create(&Repo{
			Name:      "finds" + name,
			NameSpace: "gitea",
			URL:       urls[i],
		})
	}
	store := &RepoStore{DB: db}
	repositories, err := store.Finds(urls...)
	if err != nil {
		t.Error(err)
		return
	}
	for i, repo := range repositories {
		if repo.Name != "finds"+names[i] {
			t.Fail()
		}
		if repo.ID <= 0 {
			t.Fail()
		}
	}
}

func TestRepoAssociation(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	session := db.Session()
	user := &User{
		Name:  "user",
		Email: "associate@email",
	}
	repo := &Repo{
		Name:      "repo",
		NameSpace: "associate",
		SCM:       string(core.Github),
		URL:       "http://associate/repo",
	}
	if err := session.Create(&user).Error; err != nil {
		t.Error(err)
		return
	}

	repoStore := &RepoStore{
		DB: db,
	}

	if err := repoStore.Create(repo.ToCoreRepo(), user.toCoreUser()); err != nil {
		t.Error(err)
	}

	u, err := repoStore.Creator(repo.ToCoreRepo())
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(u, user.toCoreUser()) {
		t.Fail()
	}
}

func TestRepoSetting(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()

	store := &RepoStore{DB: db}

	repo := &core.Repo{
		ID: 1234,
	}
	setting := &core.RepoSetting{
		Filters: []string{"a", "b", "c"},
	}

	if err := store.UpdateSetting(repo, setting); err != nil {
		t.Error(err)
		return
	}

	target, err := store.Setting(repo)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(target.Filters, setting.Filters) {
		t.Fail()
	}

	setting.Filters = []string{"a"}

	if err := store.UpdateSetting(repo, setting); err != nil {
		t.Error(err)
		return
	}

	target, err = store.Setting(repo)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(target.Filters, setting.Filters) {
		t.Fail()
	}

}
