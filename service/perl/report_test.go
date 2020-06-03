package perl

import (
	"context"
	"os"
	"testing"
)

func TestPerlCoverageReportService(t *testing.T) {
	f, err := os.Open("/home/coder/project/perl_test/cover_db.zip")
	if err != nil {
		t.Error(err)
		return
	}
	s, err := NewPerlCoverageReportService(f)
	if err != nil {
		t.Error(err)
		return
	}
	if len(s.db.Runs) <= 0 {
		t.Fail()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	report, err := s.Report(ctx)
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
