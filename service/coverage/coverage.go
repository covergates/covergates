package coverage

import (
	"errors"
	"io"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/service/perl"
)

var errReportTypeNotSupport = errors.New("Report type not support")

func IsReportTypeNotSupportError(err error) bool {
	if err == errReportTypeNotSupport {
		return true
	}
	return false
}

// TODO: Change the report service to do nothing constructor
func NewCoverageReportService(t core.ReportType, data io.Reader) (core.CoverageReportService, error) {
	switch t {
	case core.ReportPerl:
		return perl.NewPerlCoverageReportService(data)
	default:
		return nil, errReportTypeNotSupport
	}
}
