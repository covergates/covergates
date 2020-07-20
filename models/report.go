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
	FileData []byte
	Type     string
	ReportID string `gorm:"unique_index:report_id_commit"`
	Branch   string `gorm:"index"`
	Tag      string `gorm:"index"`
	Commit   string `gorm:"unique_index:report_id_commit"`
}

// ReportStore reports in storage
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

// Find report with the input seed. No-empty filed will use as where condition
func (store *ReportStore) Find(r *core.Report) (*core.Report, error) {
	session := store.DB.Session()
	rst := &Report{}
	session = session.Order("created_at desc").First(rst, query(r))
	if err := session.Error; err != nil {
		return nil, err
	}
	return rst.ToCoreReport(), nil
}

// Finds all report with given seed
func (store *ReportStore) Finds(r *core.Report) ([]*core.Report, error) {
	session := store.DB.Session()
	var rst []*Report
	if err := session.Find(&rst, query(r)).Error; err != nil {
		return nil, err
	}
	reports := make([]*core.Report, len(rst))
	for i, report := range rst {
		reports[i] = report.ToCoreReport()
	}
	return reports, nil
}

// CoverageReport un-marshal the coverage data
func (r *Report) CoverageReport() (*core.CoverageReport, error) {
	cover := &core.CoverageReport{}
	err := json.Unmarshal(r.Data, cover)
	return cover, err
}

// Files of the report
func (r *Report) Files() ([]string, error) {
	var files []string
	err := json.Unmarshal(r.FileData, &files)
	return files, err
}

// ToCoreReport object
func (r *Report) ToCoreReport() *core.Report {
	coverage, err := r.CoverageReport()
	if err != nil {
		coverage = &core.CoverageReport{}
	}
	files, err := r.Files()
	if err != nil {
		files = []string{}
	}
	report := &core.Report{
		Coverage:  coverage,
		Branch:    r.Branch,
		Commit:    r.Commit,
		ReportID:  r.ReportID,
		Tag:       r.Tag,
		Files:     files,
		Type:      core.ReportType(r.Type),
		CreatedAt: r.CreatedAt,
	}
	return report
}

func query(r *core.Report) *Report {
	return &Report{
		Branch:   r.Branch,
		Commit:   r.Commit,
		ReportID: r.ReportID,
		Tag:      r.Tag,
		Type:     string(r.Type),
	}
}

func copyReport(dst *Report, src *core.Report) {
	data, err := json.Marshal(src.Coverage)
	if err != nil {
		data = []byte{}
	}
	files, err := json.Marshal(src.Files)
	if err != nil {
		files = []byte{}
	}
	dst.Branch = src.Branch
	dst.Tag = src.Tag
	dst.Commit = src.Commit
	dst.ReportID = src.ReportID
	dst.Type = string(src.Type)
	dst.Data = data
	dst.FileData = files
}
