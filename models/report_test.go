package models

import (
	"encoding/json"
	"testing"
)

type MockCoverReport struct {
	Name string
}

func TestUploadReport(t *testing.T) {
	session := db.New()
	m := &Report{
		ReportID: "1234",
		Commit:   "1234",
	}
	if err := UploadReport(m); err != nil {
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
	if err := UploadReport(m); err != nil {
		t.Error(err)
		return
	}
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
