package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
)

type MockCoverReport struct {
	Name string
}

func TestReportStoreUpload(t *testing.T) {
	ctrl, service := getDatabaseService(t)
	defer ctrl.Finish()
	store := &ReportStore{
		DB: service,
	}
	m := &core.Report{
		ReportID: "1234",
		Commit:   "1234",
	}
	if err := store.Upload(m); err != nil {
		t.Error(err)
		return
	}
	r := &core.CoverageReport{
		Files: []*core.File{
			{
				Name: "mock",
			},
		},
	}
	m.Coverage = r
	if err := store.Upload(m); err != nil {
		t.Error(err)
		return
	}
	session := store.DB.Session()
	first := &Report{}
	session = session.Where(&Report{ReportID: "1234", Commit: "1234"}).First(first)
	if err := session.Error; err != nil {
		t.Error(err)
		return
	}
	r = &core.CoverageReport{}
	if err := json.Unmarshal(first.Data, r); err != nil {
		t.Error(err)
		return
	}
	if len(r.Files) <= 0 || r.Files[0].Name != "mock" {
		t.Fail()
	}
	var count int
	session = session.Where(&Report{ReportID: "1234", Commit: "1234"}).Count(&count)
	if count != 1 {
		t.Fail()
	}
}

func TestReportFind(t *testing.T) {
	ctrl, service := getDatabaseService(t)
	defer ctrl.Finish()
	store := &ReportStore{DB: service}
	session := service.Session()
	session.Create(&Report{
		Branch:   "master",
		ReportID: "report1",
		Commit:   "commit1",
	})

	session.Create(&Report{
		Branch:   "master",
		ReportID: "report1",
		Commit:   "commit2",
	})

	report := &core.Report{
		Branch:   "master",
		ReportID: "report1",
	}

	rst, err := store.Find(report)
	if err != nil {
		t.Error(err)
	}
	if rst.Commit != "commit2" {
		t.Fail()
	}
}

func TestReportUploadFiles(t *testing.T) {
	ctrl, service := getDatabaseService(t)
	defer ctrl.Finish()
	store := &ReportStore{
		DB: service,
	}
	m := &core.Report{
		ReportID: "test_upload_files",
		Commit:   "test_upload_files",
		Files:    []string{"a", "b", "c"},
	}
	if err := store.Upload(m); err != nil {
		t.Error(err)
		return
	}
	report, err := store.Find(&core.Report{
		ReportID: m.ReportID,
		Commit:   m.Commit,
	})
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(report.Files, m.Files) {
		t.Fail()
	}
}
