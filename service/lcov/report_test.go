package lcov

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/google/go-cmp/cmp"
)

type expectation struct {
	Name     string
	Lines    int
	Coverage float64
}

func (e *expectation) RoundCoverage() float64 {
	return math.Floor(e.Coverage*100) / 100
}

func testExpectation(t *testing.T, expect, result *expectation) {
	expect.Coverage = expect.RoundCoverage()
	result.Coverage = result.RoundCoverage()
	if diff := cmp.Diff(expect, result); diff != "" {
		t.Fatal(diff)
	}
}

func toExpectation(file *core.File) *expectation {
	return &expectation{
		Name:     file.Name,
		Lines:    len(file.StatementHits),
		Coverage: file.StatementCoverage,
	}
}

func mapping(files []*core.File) map[string]*core.File {
	m := make(map[string]*core.File)
	for _, file := range files {
		m[file.Name] = file
	}
	return m
}

func TestJestReport(t *testing.T) {
	s := &CoverageService{}
	file, err := os.Open(filepath.Join("testdata", "lcov.info"))
	if err != nil {
		t.Fatal(err)
	}
	report, err := s.Report(context.Background(), file)
	if err != nil {
		t.Fatal(err)
	}

	expects := []*expectation{
		{
			Name:     "/home/blueworrybear/projects/covergates/web/src/server.ts",
			Coverage: 25.0 / 47.0,
			Lines:    47,
		},
		{
			Name:     "/home/blueworrybear/projects/covergates/web/src/components/RepoListItem.vue",
			Coverage: 38.0 / 39.0,
			Lines:    39,
		},
		{
			Name:     "/home/blueworrybear/projects/covergates/web/src/store/modules/repository/mutations.ts",
			Coverage: 0,
			Lines:    10,
		},
	}
	m := mapping(report.Files)
	for _, expect := range expects {
		result, ok := m[expect.Name]
		if !ok {
			result = &core.File{}
		}
		testExpectation(t, expect, toExpectation(result))
	}
}

func TestCppReport(t *testing.T) {
	s := &CoverageService{}
	file, err := os.Open(filepath.Join("testdata", "main_coverage.info"))
	if err != nil {
		t.Fatal(err)
	}
	report, err := s.Report(context.Background(), file)
	if err != nil {
		t.Fatal(err)
	}
	expects := []*expectation{
		{
			Name:     "/home/blueworrybear/projects/test/googletest/mytest/main.cpp",
			Coverage: 1,
			Lines:    6,
		},
	}
	m := mapping(report.Files)
	for _, expect := range expects {
		result, ok := m[expect.Name]
		if !ok {
			result = &core.File{}
		}
		testExpectation(t, expect, toExpectation(result))
	}
}

func TestFind(t *testing.T) {
	s := &CoverageService{}
	p, err := s.Find(context.Background(), "testdata")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(p, "testdata/lcov.info"); diff != "" {
		t.Fatal(diff)
	}
	p, err = s.Find(context.Background(), "testdata/main_coverage.info")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(p, "testdata/main_coverage.info"); diff != "" {
		t.Fatal(diff)
	}
}
