package core

import (
	"context"
	"io"
)

//go:generate mockgen -package mock -destination ../mock/report_mock.go . ReportStore,CoverageService

// Report defined the code report structure
type Report struct {
	Coverage *CoverageReport
	Type     ReportType
	ReportID string
	Branch   string
	Tag      string
	Commit   string
}

// CoverageReport defined the code coverage report
type CoverageReport struct {
	Files             []*File
	StatementCoverage float64
}

// CoverageService provides CoverReport
type CoverageService interface {
	Report(ctx context.Context, t ReportType, r io.Reader) (*CoverageReport, error)
}

// ReportStore the report in storage
type ReportStore interface {
	Upload(r *Report) error
	Find(reportID, commit string) (*Report, error)
}
