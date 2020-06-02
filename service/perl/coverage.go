package perl

import (
	"context"
	"io"

	"github.com/code-devel-cover/CodeCover/core"
)

type PerlCoverageReportService struct {
}

func NewPerlCoverageReportService(data io.Reader) *PerlCoverageReportService {
	return nil
}

func (r *PerlCoverageReportService) Report(ctx context.Context) core.CoverageReport {
	return nil
}
