package models

import (
	"encoding/json"
	"errors"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/jinzhu/gorm"
)

var errReportID = errors.New("Error Report ID")

// Report holds the test coverage report
type Report struct {
	gorm.Model
	Data     []byte
	Type     string
	ReportID string `gorm:"unique_index:report_id_commit"`
	Branch   string `gorm:"index"`
	Tag      string `gorm:"index"`
	Commit   string `gorm:"unique_index:report_id_commit"`
}

// UploadReport create a report to database
// If the report id and commit is already existed in the table,
// the report will be updated instead.
func UploadReport(r *Report) error {
	if r.ReportID == "" || r.Commit == "" {
		return errReportID
	}
	session := db.New()
	return session.Save(r).Error
}

func FindReport(r *Report) error {
	if r.ReportID == "" || r.Commit == "" {
		return errReportID
	}
	session := db.New()
	session = session.First(r, &Report{
		ReportID: r.ReportID,
		Commit:   r.Commit,
	})
	return session.Error
}

func (r *Report) CoverageReport() (*core.CoverageReport, error) {
	cover := &core.CoverageReport{}
	err := json.Unmarshal(r.Data, cover)
	return cover, err
}
