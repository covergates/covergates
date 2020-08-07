package core

import (
	"context"
	"io"
	"time"
)

//go:generate mockgen -package mock -destination ../mock/report_mock.go . ReportStore,CoverageService

// FileNameFilters is a list of regular expression to trim file name
type FileNameFilters []string

// Report defined the code report structure
type Report struct {
	Files     []string          `json:"files"`
	Coverages []*CoverageReport `json:"coverages"`
	ReportID  string            `json:"reportID"`
	Commit    string            `json:"commit"`
	Reference string            `json:"reference"`
	CreatedAt time.Time         `json:"createdAt"`
}

// ReportComment in the pull request
type ReportComment struct {
	Number  int
	Comment int
}

// CoverageReport defined the code coverage report
type CoverageReport struct {
	Files             []*File    `json:"files"`
	Type              ReportType `json:"type"`
	StatementCoverage float64    `json:"statementCoverage"`
}

// CoverageReportDiff defines the difference between coverage reports
type CoverageReportDiff struct {
	StatementCoverageDiff float64
	Files                 []*FileDiff
}

// CoverageService provides CoverReport
type CoverageService interface {
	Report(ctx context.Context, t ReportType, r io.Reader) (*CoverageReport, error)
	// Find coverage report from the given path.
	Find(ctx context.Context, t ReportType, path string) (string, error)
	Open(ctx context.Context, t ReportType, path string) (io.Reader, error)
	// TrimFileNames in the coverage report
	TrimFileNames(ctx context.Context, report *CoverageReport, filters FileNameFilters) error
}

// ReportStore the report in storage
type ReportStore interface {
	Upload(r *Report) error
	Find(r *Report) (*Report, error)
	Finds(r *Report) ([]*Report, error)
	// List reports with reference (commit, branch or tag)
	List(reportID, ref string) ([]*Report, error)
	CreateComment(r *Report, comment *ReportComment) error
	FindComment(r *Report, number int) (*ReportComment, error)
}

// ReportService provides reports operations
type ReportService interface {
	DiffReports(source, target *Report) (*CoverageReportDiff, error)
	MarkdownReport(source, target *Report) (io.Reader, error)
	MergeReport(from, to *Report, changes []*FileChange) (*Report, error)
}

// StatementCoverage of the report
func (report *Report) StatementCoverage() float64 {
	if len(report.Coverages) <= 0 {
		return 0.0
	}
	sum := 0.0
	for _, coverage := range report.Coverages {
		sum += coverage.ComputeStatementCoverage()
	}
	return sum / float64(len(report.Coverages))
}

// Find coverage report of given type
func (report *Report) Find(t ReportType) (*CoverageReport, bool) {
	for _, coverage := range report.Coverages {
		if coverage.Type == t {
			return coverage, true
		}
	}
	return nil, false
}

// ComputeStatementCoverage of the report
func (cov *CoverageReport) ComputeStatementCoverage() float64 {
	if len(cov.Files) <= 0 {
		return 0
	}
	sum := 0.0
	for _, file := range cov.Files {
		sum += file.StatementCoverage
	}
	return sum / float64(len(cov.Files))
}
