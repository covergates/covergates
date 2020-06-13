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

type ReportStore struct {
	DB core.DatabaseService
}

// Upload create a report to database
// If the report id and commit is already existed in the table,
// the report will be updated instead.
func (store *ReportStore) Upload(r *core.Report) error {
	if r.ReportID == "" || r.Commit == "" {
		return errReportID
	}
	session := store.DB.Session()
	report := &Report{}
	session.FirstOrCreate(report, &Report{
		ReportID: r.ReportID,
		Commit:   r.Commit,
	})
	copyReport(report, r)
	return session.Save(report).Error
}

func (store *ReportStore) Find(reportID, commit string) (*core.Report, error) {
	if reportID == "" || commit == "" {
		return nil, errReportID
	}
	r := &Report{}
	session := store.DB.Session()
	session = session.First(r, &Report{
		ReportID: reportID,
		Commit:   commit,
	})
	if err := session.Error; err != nil {
		return nil, err
	}
	return r.ToCoreReport(), nil
}

func (r *Report) CoverageReport() (*core.CoverageReport, error) {
	cover := &core.CoverageReport{}
	err := json.Unmarshal(r.Data, cover)
	return cover, err
}

func (r *Report) ToCoreReport() *core.Report {
	coverage, err := r.CoverageReport()
	if err != nil {
		coverage = &core.CoverageReport{}
	}
	report := &core.Report{
		Coverage: coverage,
		Branch:   r.Branch,
		Commit:   r.Commit,
		ReportID: r.ReportID,
		Tag:      r.Tag,
		Type:     core.ReportType(r.Type),
	}
	return report
}

func copyReport(dst *Report, src *core.Report) {
	data, err := json.Marshal(src.Coverage)
	if err != nil {
		data = []byte{}
	}
	dst.Branch = src.Branch
	dst.Tag = src.Tag
	dst.Commit = src.Commit
	dst.ReportID = src.ReportID
	dst.Type = string(src.Type)
	dst.Data = data
}
