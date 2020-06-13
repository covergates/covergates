package perl

import (
	"context"
	"os"
	"testing"
)

func TestPerlCoverageService(t *testing.T) {
	f, err := os.Open("../../cover_db.zip")
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
