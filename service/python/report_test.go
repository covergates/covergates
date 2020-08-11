package python_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/service/python"
	"github.com/google/go-cmp/cmp"
)

func TestReport(t *testing.T) {
	service := &python.CoverageService{}
	file, err := os.Open(filepath.Join("testdata", "coverage.xml"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	report, err := service.Report(context.Background(), file)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Files) <= 0 || report.StatementCoverage <= 0 {
		t.Fatal()
	}
	m := make(map[string]bool)
	for _, file := range report.Files {
		m[file.Name] = true
	}
	targets := []string{
		"/path/to/python/project/__init__.py",
		"/path/to/python/project/apis/qa.py",
	}
	for _, target := range targets {
		if _, ok := m[target]; !ok {
			t.Fatalf("cannot find %s", target)
		}
	}

}

func TestFind(t *testing.T) {
	service := &python.CoverageService{}
	report, err := service.Find(context.Background(), "testdata")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(report) != "coverage.xml" {
		t.Fatal()
	}
}

func TestReportSimple(t *testing.T) {
	service := &python.CoverageService{}
	file, err := os.Open(filepath.Join("testdata", "coverage_simple.xml"))
	if err != nil {
		t.Fatal(err)
	}
	report, err := service.Report(context.Background(), file)
	if err != nil {
		t.Fatal(err)
	}
	root := "/home/blueworrybear/projects/test/pytest"
	expectFiles := map[string]*core.File{
		filepath.Join(root, "m1/tests/test_f1.py"): {
			Name:              filepath.Join(root, "m1/tests/test_f1.py"),
			StatementCoverage: float64(1),
			StatementHits: []*core.StatementHit{
				{
					Hits:       1,
					LineNumber: 1,
				},
				{
					Hits:       1,
					LineNumber: 4,
				},
				{
					Hits:       1,
					LineNumber: 5,
				},
			},
		},
		filepath.Join(root, "m2/f2.py"): {
			Name:              filepath.Join(root, "m2/f2.py"),
			StatementCoverage: float64(1),
			StatementHits: []*core.StatementHit{
				{
					Hits:       1,
					LineNumber: 2,
				},
				{
					Hits:       1,
					LineNumber: 3,
				},
			},
		},
	}

	if report.Type != core.ReportPython {
		t.Fatal()
	}

	m := make(map[string]*core.File)
	for _, file := range report.Files {
		m[file.Name] = file
	}

	for name := range expectFiles {
		expect := expectFiles[name]
		file, ok := m[name]
		if !ok {
			t.Fatalf("cannot find %s", name)
		}
		if diff := cmp.Diff(expect, file); diff != "" {
			t.Fatal(diff)
		}
	}

}
