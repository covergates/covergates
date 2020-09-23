package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUpload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if diff := cmp.Diff("Bearer token", token); diff != "" {
			t.Fatal(diff)
		}
		if diff := cmp.Diff("/reports/123", r.URL.String()); diff != "" {
			t.Fatal(diff)
		}
	}))
	defer ts.Close()
	os.Setenv("API_URL", ts.URL)
	os.Setenv("REPORT_ID", "123")
	os.Setenv("GATES_TOKEN", "token")
	args := []string{"", "upload", "--type", "go", "./testdata/coverage.out"}
	if err := app.Run(args); err != nil {
		t.Fatal(err)
	}

}
