package models

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/google/go-cmp/cmp"
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
		Type:     core.ReportPerl,
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

func TestReportFinds(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := &ReportStore{DB: db}
	// TODO: Add more report types
	reportID := "find_reports"
	commit := "abcdefgcomment"
	now := time.Now()
	reports := []*core.Report{
		{
			ReportID:  reportID,
			Coverage:  &core.CoverageReport{},
			Commit:    commit,
			Type:      core.ReportPerl,
			CreatedAt: now,
		},
	}
	for _, report := range reports {
		if err := store.Upload(report); err != nil {
			t.Fatal(err)
		}
	}

	found, err := store.Finds(&core.Report{ReportID: reportID, Commit: commit})
	if err != nil {
		t.Fatal(err)
	}
	for _, report := range found {
		report.CreatedAt = now
	}
	if diff := cmp.Diff(found, reports); diff != "" {
		t.Fatal(diff)
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
		Type:     core.ReportPerl,
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

func TestReportComment(t *testing.T) {
	ctrl, service := getDatabaseService(t)
	defer ctrl.Finish()
	store := &ReportStore{
		DB: service,
	}

	report := &core.Report{
		ReportID: "ABCD",
	}

	if err := store.CreateComment(report, &core.ReportComment{}); err == nil {
		t.Fail()
	}

	if err := store.CreateComment(report, &core.ReportComment{Comment: 1, Number: 1}); err != nil {
		t.Fatal(err)
	}
	comment, err := store.FindComment(report, 1)
	if err != nil {
		t.Fatal(err)
	}
	if comment.Comment != 1 {
		t.Fail()
	}
	if err := store.CreateComment(report, &core.ReportComment{Comment: 2, Number: 1}); err != nil {
		t.Fatal(err)
	}
	comment, err = store.FindComment(report, 1)
	if err != nil {
		t.Fatal(err)
	}
	if comment.Comment != 2 {
		t.Fail()
	}
	if _, err := store.FindComment(report, 123); err == nil {
		t.Fail()
	}
}
