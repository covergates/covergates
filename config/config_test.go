package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func setEnv(vars map[string]string) {
	for env := range vars {
		os.Setenv(env, vars[env])
	}
}

func unsetEnv(vars map[string]string) {
	for env := range vars {
		os.Unsetenv(env)
	}
}

func TestEnvironConfig(t *testing.T) {

	cfg, err := Environ()

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(cfg.Server.Addr, "http://localhost:8080"); diff != "" {
		t.Log(diff)
		t.Fail()
	}
	vars := map[string]string{
		"GATES_SERVER_ADDR": "http://localhost:3000",
	}
	setEnv(vars)

	if cfg, err = Environ(); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(vars["GATES_SERVER_ADDR"], cfg.Server.Addr); diff != "" {
		t.Log(diff)
		t.Fail()
	}
	unsetEnv(vars)

	vars = map[string]string{
		"GATES_GITEA_SERVER":        "http://localhost:3000",
		"GATES_GITEA_CLIENT_ID":     "c8c6a2cc-f948-475c-8663-f420c8fc15ab",
		"GATES_GITEA_CLIENT_SECRET": "J8YYirhYOZY9a9RepaoORN-8EFcSO-sbwjSGvGo4NwE=",
		"GATES_GITEA_SCOPE":         "repo,repo:status",
		"GATES_GITEA_SKIP_VERITY":   "true",
	}

	setEnv(vars)

	if cfg, err = Environ(); err != nil {
		t.Fatal(err)
	}

	result := map[string]string{
		"GATES_GITEA_SERVER":        cfg.Gitea.Server,
		"GATES_GITEA_CLIENT_ID":     cfg.Gitea.ClientID,
		"GATES_GITEA_CLIENT_SECRET": cfg.Gitea.ClientSecret,
		"GATES_SERVER_BASE":         cfg.Server.Base,
	}

	expect := map[string]string{
		"GATES_GITEA_SERVER":        "http://localhost:3000",
		"GATES_GITEA_CLIENT_ID":     "c8c6a2cc-f948-475c-8663-f420c8fc15ab",
		"GATES_GITEA_CLIENT_SECRET": "J8YYirhYOZY9a9RepaoORN-8EFcSO-sbwjSGvGo4NwE=",
		"GATES_SERVER_BASE":         "",
	}

	if diff := cmp.Diff(result, expect); diff != "" {
		t.Log(diff)
		t.Fail()
	}

	if diff := cmp.Diff(cfg.Gitea.Scope, []string{"repo", "repo:status"}); diff != "" {
		t.Log(diff)
		t.Fail()
	}

	if !cfg.Gitea.SkipVerity {
		t.Fail()
	}

}

func TestConfig(t *testing.T) {
	vars := map[string]string{
		"GATES_SERVER_ADDR": "http://localhost:3000",
		"GATES_SERVER_BASE": "/gates",
	}

	setEnv(vars)
	defer unsetEnv(vars)

	cfg, err := Environ()
	if err != nil {
		t.Fatal(err)
	}

	result := []interface{}{cfg.Server.Port(), cfg.Server.BaseURL()}
	expect := []interface{}{"3000", "/gates"}

	if diff := cmp.Diff(result, expect); diff != "" {
		t.Log(diff)
		t.Fail()
	}
}
