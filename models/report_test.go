package models

import (
	"encoding/json"
	"testing"
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
	m := &Report{
		ReportID: "1234",
		Commit:   "1234",
	}
	if err := store.UploadReport(m); err != nil {
		t.Error(err)
		return
	}
	r := &MockCoverReport{
		Name: "mock",
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
		return
	}
	m.Data = data
	if err := store.UploadReport(m); err != nil {
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
	r = &MockCoverReport{}
	if err := json.Unmarshal(first.Data, r); err != nil {
		t.Error(err)
		return
	}
	if r.Name != "mock" {
		t.Fail()
	}
	var count int
	session = session.Where(&Report{ReportID: "1234", Commit: "1234"}).Count(&count)
	if count != 1 {
		t.Fail()
	}
}
