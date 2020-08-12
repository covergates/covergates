package ruby

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/google/go-cmp/cmp"
)

type expect struct {
	Name     string
	Lines    int
	Coverage float64
}

func TestReport(t *testing.T) {
	s := &CoverageService{}
	file, err := os.Open(filepath.Join("testdata", ".resultset.json"))
	if err != nil {
		t.Fatal(err)
	}
	r, err := s.Report(context.Background(), file)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Files) <= 0 {
		t.Fatal()
	}
	m := make(map[string]*core.File)
	for _, file := range r.Files {
		m[file.Name] = file
	}

	expects := []*expect{
		{
			Name:     "/home/blueworrybear/projects/discourse/lib/discourse_event.rb",
			Lines:    12,
			Coverage: 0.83,
		},
		{
			Name:     "/home/blueworrybear/projects/discourse/spec/fabricators/watched_word_fabricator.rb",
			Lines:    3,
			Coverage: 0.33,
		},
		{
			Name:     "/home/blueworrybear/projects/discourse/config/environment.rb",
			Lines:    4,
			Coverage: 0.75,
		},
	}

	for _, e := range expects {
		result := &expect{}
		if file, ok := m[e.Name]; ok {
			result.Name = file.Name
			result.Lines = len(file.StatementHits)
			result.Coverage = math.Floor(file.StatementCoverage*100) / 100
		}
		if diff := cmp.Diff(e, result); diff != "" {
			t.Fatal(diff)
		}
	}
}

func TestFind(t *testing.T) {
	s := &CoverageService{}
	if _, err := s.Find(context.Background(), "testdata"); err != nil {
		t.Fatal()
	}
	if _, err := s.Find(context.Background(), filepath.Join("testdata", ".resultset.json")); err != nil {
		t.Fatal()
	}
}
