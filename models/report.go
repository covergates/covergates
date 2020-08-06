package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/covergates/covergates/core"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var errReportFields = errors.New("Error Report Fields")

// Report holds the test coverage report
type Report struct {
	gorm.Model
	Data       []byte
	FileData   []byte
	Type       string       `gorm:"unique_index:report_record"`
	ReportID   string       `gorm:"unique_index:report_record"`
	References []*Reference `gorm:"many2many:report_reference"`
	Commit     string       `gorm:"unique_index:report_record"`
}

// Reference of Report, such as branch or tag name
type Reference struct {
	gorm.Model
	ReportID string    `gorm:"unique_index:reference_record"`
	Name     string    `gorm:"unique_index:reference_record"`
	Reports  []*Report `gorm:"many2many:report_reference"`
}

// ReportComment defines summary report comment in the pull request
type ReportComment struct {
	gorm.Model
	ReportID string `gorm:"unique_index:report_comment_number"`
	// Number is the PR number
	Number  int `gorm:"unique_index:report_comment_number"`
	Comment int
}

// ReportStore reports in storage
type ReportStore struct {
	DB core.DatabaseService
}

// Upload create a report to database
// If the report id and commit is already existed in the table,
// the report will be updated instead.
func (store *ReportStore) Upload(r *core.Report) error {
	if r.ReportID == "" || r.Commit == "" || r.Type == "" {
		return errReportFields
	}
	session := store.DB.Session()
	if r.Reference != "" {
		session = session.Preload("References", "name=?", r.Reference)
	}
	report := &Report{}
	session.FirstOrCreate(report, &Report{
		ReportID: r.ReportID,
		Commit:   r.Commit,
		Type:     string(r.Type),
	})
	if len(report.References) <= 0 && r.Reference != "" {
		session := store.DB.Session()
		ref := &Reference{Name: r.Reference, ReportID: r.ReportID}
		if err := session.FirstOrCreate(ref, &ref).Error; err != nil {
			return err
		}
		if err := session.Model(&report).Association("References").Append(ref).Error; err != nil {
			return err
		}
	}
	copyReport(report, r)
	return session.Save(report).Error
}

// Find report with the input seed. No-empty filed will use as where condition
func (store *ReportStore) Find(r *core.Report) (*core.Report, error) {
	session := store.DB.Session()
	target := &Report{}
	if r.Reference == "" {
		session = session.Order("created_at desc").First(target, query(r))
		if err := session.Error; err != nil {
			return nil, err
		}
	} else {
		if r.ReportID == "" {
			log.Warning("report id should not be empty when search with reference")
			return nil, gorm.ErrRecordNotFound
		}
		ref := &Reference{ReportID: r.ReportID, Name: r.Reference}
		session = session.Preload("Reports", func(db *gorm.DB) *gorm.DB {
			return db.Where(query(r)).Order("created_at desc")
		}).First(ref, &ref)
		if err := session.Error; err != nil {
			return nil, err
		}
		if len(ref.Reports) <= 0 {
			return nil, gorm.ErrRecordNotFound
		}
		target = ref.Reports[0]
	}
	report := target.ToCoreReport()
	report.Reference = r.Reference
	return report, nil
}

// Finds all report with given seed
func (store *ReportStore) Finds(r *core.Report) ([]*core.Report, error) {
	session := store.DB.Session()
	var rst []*Report

	if r.Reference == "" {
		if err := session.Where(query(r)).Find(&rst).Error; err != nil {
			return nil, err
		}
	} else {
		if r.ReportID == "" {
			log.Warning("report id should not be empty when search with reference")
			return nil, gorm.ErrRecordNotFound
		}
		ref := &Reference{ReportID: r.ReportID, Name: r.Reference}
		session = session.Preload("Reports", func(db *gorm.DB) *gorm.DB {
			return db.Where(query(r)).Order("created_at desc")
		}).First(ref, &ref)
		if err := session.Error; err != nil {
			return nil, err
		}
		if len(ref.Reports) <= 0 {
			return nil, gorm.ErrRecordNotFound
		}
		rst = ref.Reports
	}
	reports := make([]*core.Report, len(rst))
	for i, report := range rst {
		reports[i] = report.ToCoreReport()
		reports[i].Reference = r.Reference
	}
	return reports, nil
}

// CreateComment of the report summary
func (store *ReportStore) CreateComment(r *core.Report, comment *core.ReportComment) error {
	if comment.Comment <= 0 || comment.Number <= 0 {
		return fmt.Errorf("invalid comment")
	}
	session := store.DB.Session()
	condition := &ReportComment{ReportID: r.ReportID, Number: comment.Number}
	c := &ReportComment{}
	if err := session.Where(condition).FirstOrCreate(c).Error; err != nil {
		return err
	}
	c.Comment = comment.Comment
	return session.Save(c).Error
}

// FindComment summary of given PR number
func (store *ReportStore) FindComment(r *core.Report, number int) (*core.ReportComment, error) {
	session := store.DB.Session()
	condition := &ReportComment{ReportID: r.ReportID, Number: number}
	comment := &ReportComment{}
	if err := session.First(comment, condition).Error; err != nil {
		return nil, err
	}
	return &core.ReportComment{
		Comment: comment.Comment,
		Number:  comment.Number,
	}, nil
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
		Commit:    r.Commit,
		ReportID:  r.ReportID,
		Files:     files,
		Type:      core.ReportType(r.Type),
		CreatedAt: r.CreatedAt,
	}
	return report
}

func query(r *core.Report) *Report {
	return &Report{
		Commit:   r.Commit,
		ReportID: r.ReportID,
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
	dst.Commit = src.Commit
	dst.ReportID = src.ReportID
	dst.Type = string(src.Type)
	dst.Data = data
	dst.FileData = files
}
