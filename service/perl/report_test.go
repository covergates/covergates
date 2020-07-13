package perl

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPerlCoverageService(t *testing.T) {
	file := filepath.Join("testdata", "cover_db.zip")
	f, err := os.Open(file)
	if err != nil {
		t.Error(err)
		return
	}
	s := &CoverageService{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	report, err := s.Report(ctx, f)
	if err != nil {
		t.Error(err)
		return
	}
	if len(report.Files) <= 0 {
		t.Fail()
		return
	}
	if report.StatementCoverage <= 0 {
		t.Fail()
		return
	}
}

func TestPerlCoverageWithStringStatementHits(t *testing.T) {
	file := filepath.Join("testdata", "JSON_db.zip")
	f, err := os.Open(file)
	if err != nil {
		t.Error(err)
		return
	}
	s := &CoverageService{}
	report, err := s.Report(context.Background(), f)
	if err != nil {
		t.Error(err)
		return
	}
	if len(report.Files) <= 0 {
		t.Fail()
		return
	}
	if report.StatementCoverage <= 0 {
		t.Fail()
		return
	}
}
