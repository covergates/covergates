package python_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/covergates/covergates/service/python"
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
