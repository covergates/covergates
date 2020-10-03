package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/covergates/covergates/core"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var errReportFields = errors.New("Error Report Fields")

type reportList []*Report

// Report holds the report
type Report struct {
	gorm.Model
	FileData   []byte
	ReportID   string `gorm:"size:256;uniqueIndex:report_record"`
	Coverages  []*Coverage
	References []*Reference `gorm:"many2many:report_reference"`
	Commit     string       `gorm:"size:256;uniqueIndex:report_record"`
}

// Coverage defines test coverage report
type Coverage struct {
	gorm.Model
	Data     []byte
	Type     string
	ReportID uint
}

// Reference of Report, such as branch or tag name
type Reference struct {
	gorm.Model
	ReportID string    `gorm:"size:256;uniqueIndex:reference_record"`
	Name     string    `gorm:"size:256;uniqueIndex:reference_record"`
	Reports  []*Report `gorm:"many2many:report_reference"`
}

// ReportComment defines summary report comment in the pull request
type ReportComment struct {
	gorm.Model
	ReportID string `gorm:"size:256;uniqueIndex:report_comment_number"`
	// Number is the PR number
	Number  int `gorm:"uniqueIndex:report_comment_number"`
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
	if r.ReportID == "" || r.Commit == "" {
		return errReportFields
	}
	session := store.DB.Session()
	if r.Reference != "" {
		session = session.Preload("References", "name=?", r.Reference)
	}
	report := &Report{}
	session.Preload("Coverages").FirstOrCreate(report, &Report{
		ReportID: r.ReportID,
		Commit:   r.Commit,
	})
	if len(report.References) <= 0 && r.Reference != "" {
		if err := store.appendReference(report, r.Reference); err != nil {
			return err
		}
	}
	for _, coverage := range r.Coverages {
		if err := store.updateCoverage(report, coverage); err != nil {
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
		session = session.Preload("Coverages").Order("created_at desc").First(target, query(r))
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
			return db.Where(query(r)).Order("created_at desc").Limit(50)
		}).Preload("Reports.Coverages").First(ref, ref)
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
		if err := session.Preload(
			"Coverages",
		).Where(query(r)).Order(
			"created_at desc",
		).Limit(100).Find(&rst).Error; err != nil {
			return nil, err
		}
	} else {
		if r.ReportID == "" {
			log.Warning("report id should not be empty when search with reference")
			return nil, gorm.ErrRecordNotFound
		}
		ref := &Reference{ReportID: r.ReportID, Name: r.Reference}
		session = session.Preload("Reports", func(db *gorm.DB) *gorm.DB {
			return db.Where(query(r)).Order("created_at desc").Limit(100)
		}).Preload("Reports.Coverages").First(ref, ref)
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

// List reports with reference
//
// reference (ref) could be commit SHA, branch or tag name.
// The files and data field will be remove from result to reduce memory usage.
func (store *ReportStore) List(reportID, ref string) ([]*core.Report, error) {
	session := store.DB.Session()
	var reports reportList
	condition := &Report{ReportID: reportID, Commit: ref}
	err := session.Preload("Coverages").Where(condition).Order(
		"created_at desc",
	).Limit(100).Find(&reports).Error
	if err == nil && len(reports) > 0 {
		return reports.ToCoreReports(""), nil
	}
	reference := &Reference{ReportID: reportID, Name: ref}
	session = store.DB.Session().Preload("Reports", func(db *gorm.DB) *gorm.DB {
		return db.Order(
			"created_at desc",
		).Limit(200)
	}).Preload("Reports.Coverages").First(reference, reference)
	if err := session.Error; err != nil {
		return nil, err
	}
	if len(reference.Reports) <= 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return reportList(reference.Reports).ToCoreReports(ref), nil
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

func (store *ReportStore) updateCoverage(r *Report, cov *core.CoverageReport) error {
	c, ok := r.find(cov.Type)
	if err := copyCoverage(c, cov); err != nil {
		return err
	}
	if !ok {
		r.Coverages = append(r.Coverages, c)
	} else if c.ID > 0 {
		store.DB.Session().Save(c)
	}
	return nil
}

func (store *ReportStore) appendReference(r *Report, name string) error {
	if r.ReportID == "" {
		return errReportFields
	}
	session := store.DB.Session()
	ref := &Reference{Name: name, ReportID: r.ReportID}
	if err := session.FirstOrCreate(ref, ref).Error; err != nil {
		return err
	}
	return session.Model(r).Association("References").Append(ref)
}

// Files of the report
func (r *Report) Files() ([]string, error) {
	var files []string
	err := json.Unmarshal(r.FileData, &files)
	return files, err
}

// ToCoreReport object
func (r *Report) ToCoreReport() *core.Report {
	coverages := make([]*core.CoverageReport, len(r.Coverages))
	for i, coverage := range r.Coverages {
		c, _ := coverage.ToCoreCoverage()
		coverages[i] = c
	}
	files, err := r.Files()
	if err != nil {
		files = []string{}
	}
	report := &core.Report{
		Commit:    r.Commit,
		ReportID:  r.ReportID,
		Files:     files,
		Coverages: coverages,
		CreatedAt: r.CreatedAt,
	}
	return report
}

func (r *Report) find(t core.ReportType) (*Coverage, bool) {
	for _, coverage := range r.Coverages {
		if coverage.Type == string(t) {
			return coverage, true
		}
	}
	return &Coverage{}, false
}

// ToCoreCoverage unmarshal Coverage Data
func (c *Coverage) ToCoreCoverage() (*core.CoverageReport, error) {
	cover := &core.CoverageReport{}
	cover.Type = core.ReportType(c.Type)
	if err := json.Unmarshal(c.Data, cover); err != nil {
		return cover, err
	}
	cover.StatementCoverage = cover.ComputeStatementCoverage()
	return cover, nil
}

func (r reportList) ToCoreReports(ref string) []*core.Report {
	result := make([]*core.Report, len(r))
	for i, report := range r {
		report.FileData = nil
		coreReport := report.ToCoreReport()
		for _, coverage := range coreReport.Coverages {
			coverage.Files = nil
		}
		coreReport.Reference = ref
		result[i] = coreReport
	}
	return result
}

func query(r *core.Report) *Report {
	return &Report{
		Commit:   r.Commit,
		ReportID: r.ReportID,
	}
}

func copyReport(dst *Report, src *core.Report) {
	files, err := json.Marshal(src.Files)
	if err != nil {
		files = []byte{}
	}
	dst.Commit = src.Commit
	dst.ReportID = src.ReportID
	dst.FileData = files
}

func copyCoverage(dst *Coverage, src *core.CoverageReport) error {
	cov, err := json.Marshal(src)
	if err != nil {
		return err
	}
	dst.Data = cov
	dst.Type = string(src.Type)
	return nil
}
