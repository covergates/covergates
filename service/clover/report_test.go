package clover_test

import (
	"context"
	"math"
	"os"
	"path"
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/service/clover"
)

func TestReport(t *testing.T) {
	s := &clover.CoverageService{}
	file, err := os.Open(path.Join("testdata", "library-coverage.xml"))
	if err != nil {
		t.Fatal(err)
	}
	report, err := s.Report(context.Background(), file)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Files) <= 0 {
		t.Fatal()
	}
	m := make(map[string]*core.File)
	for _, file := range report.Files {
		m[file.Name] = file
	}
	coreFile, ok := m["/home/blueworrybear/projects/vallina/vanilla/applications/vanilla/tests/Utils/CommunityApiTestTrait.php"]
	if !ok {
		t.Fatal()
	}
	if coreFile.StatementCoverage != 0 {
		t.Fatal(s)
	}
	coreFile, ok = m["/home/blueworrybear/projects/vallina/vanilla/library/Garden/EventManager.php"]
	if !ok {
		t.Fatal()
	}
	if math.Round(coreFile.StatementCoverage*100)/100 != 0.18 {
		t.Fatal(coreFile.StatementCoverage)
	}
}
