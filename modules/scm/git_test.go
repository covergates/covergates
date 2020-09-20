package scm

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGiteaCommitQuery(t *testing.T) {
	u := mustGetGiteaCommitsQuery("gitea/test", "bear")
	expect := "api/v1/repos/gitea/test/commits?sha=bear"
	if diff := cmp.Diff(expect, u); diff != "" {
		t.Fatal(diff)
	}
}
