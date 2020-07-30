package scm

import (
	"testing"

	"github.com/covergates/covergates/core"
)

func TestNewReportID(t *testing.T) {
	s := &repoService{}
	reportID1 := s.NewReportID(&core.Repo{
		URL: "http://repo",
	})
	reportID2 := s.NewReportID(&core.Repo{
		URL: "http://repo",
	})
	if reportID1 == reportID2 {
		t.Logf("%s %s", reportID1, reportID2)
		t.Fail()
	}
}
