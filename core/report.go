package core

import (
	"context"
)

type CoverageReport struct {
	Files             []*File
	StatementCoverage float32
}

type CoverageReportService interface {
	Report(ctx context.Context) *CoverageReport
}
