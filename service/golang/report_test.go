package golang

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/covergates/covergates/core"
)

func TestReport(t *testing.T) {
	s := &CoverageService{}
	file := filepath.Join("testdata", "coverage.out")
	reader, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	report, err := s.Report(context.Background(), reader)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Files) != 5 {
		t.Fatal("file count is not correct")
	}
	files := make(map[string]*core.File)
	for _, file := range report.Files {
		files[file.Name] = file
	}
	if file, ok := files["github.com/covergates/covergates/models/repo.go"]; ok {
		hits := hitMap(file.StatementHits)
		expect := map[int]bool{
			145: false,
			169: true,
			178: false,
		}
		checkHits(t, hits, expect)
	} else {
		t.Fatal("repo.go not found")
	}

}

func hitMap(hits []*core.StatementHit) map[int]*core.StatementHit {
	m := make(map[int]*core.StatementHit)
	for _, hit := range hits {
		m[hit.LineNumber] = hit
	}
	return m
}

func checkHits(t *testing.T, m map[int]*core.StatementHit, expect map[int]bool) {
	for line, hitted := range expect {
		if hit, ok := m[line]; ok {
			if hit.Hits > 0 && !hitted {
				t.Fatalf("line %d is hit", line)
			} else if hit.Hits <= 0 && hitted {
				t.Fatalf("line %d is not hit", line)
			}
		} else {
			t.Fatalf("line %d has not report", line)
		}
	}
}
