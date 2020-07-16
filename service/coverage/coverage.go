package coverage

import (
	"context"
	"errors"
	"io"
	"regexp"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/service/perl"
)

var errReportTypeNotSupport = errors.New("Report type not support")

type CoverageService struct{}
type TypeCoverageService interface {
	Report(ctx context.Context, data io.Reader) (*core.CoverageReport, error)
}

func IsReportTypeNotSupportError(err error) bool {
	if err == errReportTypeNotSupport {
		return true
	}
	return false
}

func (s *CoverageService) service(t core.ReportType) (TypeCoverageService, error) {
	switch t {
	case core.ReportPerl:
		return &perl.CoverageService{}, nil
	default:
		return nil, errReportTypeNotSupport
	}
}

func (s *CoverageService) Report(ctx context.Context, t core.ReportType, r io.Reader) (*core.CoverageReport, error) {
	service, err := s.service(t)
	if err != nil {
		return nil, err
	}
	return service.Report(ctx, r)
}

func (s *CoverageService) TrimFileNames(ctx context.Context, report *core.CoverageReport, filters core.FileNameFilters) error {
	regexps := toRegexps(filters)
	for _, file := range report.Files {
		for _, regex := range regexps {
			file.Name = regex.ReplaceAllString(file.Name, "")
		}
	}
	return nil
}

func toRegexps(slice []string) []*regexp.Regexp {
	regex := make([]*regexp.Regexp, 0, len(slice))
	for _, expr := range slice {
		if r, err := regexp.Compile(expr); err == nil {
			regex = append(regex, r)
		}
	}
	return regex
}
