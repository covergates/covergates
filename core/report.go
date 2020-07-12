package core

import (
	"context"
	"io"
)

//go:generate mockgen -package mock -destination ../mock/report_mock.go . ReportStore,CoverageService

// Report defined the code report structure
type Report struct {
	Coverage *CoverageReport `json:"coverage"`
	Files    []string        `json:"files"`
	Type     ReportType      `json:"type"`
	ReportID string          `json:"reportID"`
	Branch   string          `json:"branch"`
	Tag      string          `json:"tag"`
	Commit   string          `json::"commit"`
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
	Find(r *Report) (*Report, error)
	Finds(r *Report) ([]*Report, error)
}
