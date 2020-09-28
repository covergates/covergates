package models

import (
	"reflect"
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
)

func toCoreRepoSlice(repos []*Repo) []*core.Repo {
	result := make([]*core.Repo, len(repos))
	for i, repo := range repos {
		result[i] = repo.ToCoreRepo()
	}
	return result
}

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

func TestRepoUpdateOrCreate(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := &RepoStore{DB: db}

	repo := &core.Repo{
		SCM:     core.Github,
		URL:     "http://testrepo/update/create/1",
		Private: true,
	}

	if err := store.UpdateOrCreate(repo); err != nil {
		t.Fatal(err)
	}

	result, err := store.Find(&core.Repo{URL: repo.URL})
	if err != nil {
		t.Fatal(err)
	}
	repo.ID = result.ID
	if diff := cmp.Diff(repo, result); diff != "" {
		t.Fatal(diff)
	}

	// should not update report id
	if err := store.UpdateOrCreate(&core.Repo{SCM: repo.SCM, URL: repo.URL, Private: true}); err != nil {
		t.Fatal()
	}

	result, err = store.Find(&core.Repo{URL: repo.URL})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(repo, result); diff != "" {
		t.Fatal(diff)
	}

	// should update private
	repo.Private = false
	if err := store.UpdateOrCreate(repo); err != nil {
		t.Fatal()
	}
	result, err = store.Find(&core.Repo{URL: repo.URL})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(repo, result); diff != "" {
		t.Fatal(diff)
	}
}

func TestBatchUpdateOrCreate(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := &RepoStore{DB: db}

	repos := []*core.Repo{
		{
			SCM: core.Github,
			URL: "http://test/batch/update1",
		},
		{
			SCM: core.Github,
			URL: "http://test/batch/update2",
		},
	}

	if err := store.BatchUpdateOrCreate(repos); err != nil {
		t.Fatal(err)
	}
	urls := make([]string, len(repos))
	for i, repo := range repos {
		urls[i] = repo.URL
	}

	var result []*Repo
	if err := db.Session().Where("url IN ?", urls).Find(&result).Error; err != nil {
		t.Fatal()
	}
	for _, repo := range result {
		repo.ID = 0
	}
	if diff := cmp.Diff(toCoreRepoSlice(result), repos); diff != "" {
		t.Fatal(diff)
	}

	repos = []*core.Repo{
		{
			SCM: core.Github,
			URL: "http://test/batch/update3",
		},
		{
			URL: "",
		},
	}

	if err := store.BatchUpdateOrCreate(repos); err == nil {
		t.Fatal("should encounter error")
	}

	if err := db.Session().Where(&Repo{URL: repos[0].URL}).First(&Repo{}).Error; err != gorm.ErrRecordNotFound {
		t.Fatal()
	}
}

func TestRepoAssociation(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	session := db.Session()
	user := &User{
		Login: "test_repo",
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

	if err := repoStore.Create(repo.ToCoreRepo()); err != nil {
		t.Error(err)
	}
	if err := repoStore.UpdateCreator(repo.ToCoreRepo(), user.toCoreUser()); err != nil {
		t.Fatal(err)
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

	repo1 := &core.Repo{
		ID: 1234,
	}
	setting1 := &core.RepoSetting{
		Filters: []string{"a", "b", "c"},
	}
	repo2 := &core.Repo{
		ID: 4567,
	}
	setting2 := &core.RepoSetting{
		Filters: []string{"e", "f"},
	}

	if err := store.UpdateSetting(repo1, setting1); err != nil {
		t.Error(err)
		return
	}

	if err := store.UpdateSetting(repo2, setting2); err != nil {
		t.Error(err)
		return
	}

	target, err := store.Setting(repo1)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(target.Filters, setting1.Filters) {
		t.Logf("\nExpect: %v\nGot: %v\n", setting1.Filters, target.Filters)
		t.Fail()
	}

	target2, err := store.Setting(repo2)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(target2.Filters, setting2.Filters) {
		t.Logf("\nExpect: %v\nGot: %v\n", setting2.Filters, target2.Filters)
		t.Fail()
	}

	setting1.Filters = []string{"a"}

	if err := store.UpdateSetting(repo1, setting1); err != nil {
		t.Error(err)
		return
	}

	target, err = store.Setting(repo1)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(target.Filters, setting1.Filters) {
		t.Fail()
	}
}

func TestPrivateRepository(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := &RepoStore{DB: db}

	user := &core.User{Login: "user"}

	coreRepo := &core.Repo{
		Name:      "private_repo",
		NameSpace: "gitea",
		URL:       "private_repo.com",
		SCM:       core.Gitea,
	}

	if err := store.Create(coreRepo); err != nil {
		t.Fatal(err)
	}
	if err := store.UpdateCreator(coreRepo, user); err != nil {
		t.Fatal(err)
	}

	repo, err := store.Find(&core.Repo{Name: coreRepo.Name})
	if err != nil {
		t.Fatal(err)
	}

	if repo.Private {
		t.Log("repo should be public")
		t.Fail()
	}

	coreRepo.Private = true

	if err := store.Update(coreRepo); err != nil {
		t.Fatal(err)
	}

	if repo, err = store.Find(&core.Repo{Name: coreRepo.Name}); err != nil {
		t.Fatal(err)
	}

	if !repo.Private {
		t.Log("repo should be private")
		t.Fail()
	}
}

func TestRepoHook(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := &RepoStore{DB: db}

	if err := store.UpdateHook(&core.Repo{}, &core.Hook{}); err == nil {
		t.Fatal("invalid repo should return error")
	}

	repo := &core.Repo{
		ID: uint(123),
	}

	if err := store.UpdateHook(repo, &core.Hook{}); err == nil {
		t.Fatal("invalid hook should return error")
	}

	if _, err := store.FindHook(repo); err == nil {
		t.Fail()
	}
	expectHooks := []*core.Hook{
		{
			ID: "1",
		},
		{
			ID: "2",
		},
	}
	for _, expect := range expectHooks {
		if err := store.UpdateHook(repo, expect); err != nil {
			t.Fatal(err)
		}
		hook, err := store.FindHook(repo)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(hook, expect); diff != "" {
			t.Log(diff)
			t.Fail()
		}
	}
}
