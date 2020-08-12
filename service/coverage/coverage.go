package coverage

import (
	"context"
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/service/golang"
	"github.com/covergates/covergates/service/lcov"
	"github.com/covergates/covergates/service/perl"
	"github.com/covergates/covergates/service/python"
	"github.com/covergates/covergates/service/ruby"
)

var errReportTypeNotSupport = errors.New("Report type not support")

// Service of coverage report
type Service struct{}

//TypeCoverageService defines a coverage service for a language
type TypeCoverageService interface {
	Report(ctx context.Context, data io.Reader) (*core.CoverageReport, error)
	Find(ctx context.Context, path string) (string, error)
	Open(ctx context.Context, path string) (io.Reader, error)
}

// IsReportTypeNotSupportError check
func IsReportTypeNotSupportError(err error) bool {
	if err == errReportTypeNotSupport {
		return true
	}
	return false
}

func (s *Service) service(t core.ReportType) (TypeCoverageService, error) {
	switch t {
	case core.ReportPerl:
		return &perl.CoverageService{}, nil
	case core.ReportGo:
		return &golang.CoverageService{}, nil
	case core.ReportPython:
		return &python.CoverageService{}, nil
	case core.ReportRuby:
		return &ruby.CoverageService{}, nil
	case core.ReportLCOV:
		return &lcov.CoverageService{}, nil
	default:
		return nil, errReportTypeNotSupport
	}
}

// Report coverage from reader data
func (s *Service) Report(ctx context.Context, t core.ReportType, r io.Reader) (*core.CoverageReport, error) {
	service, err := s.service(t)
	if err != nil {
		return nil, err
	}
	return service.Report(ctx, r)
}

// Find coverage report data from given path
func (s *Service) Find(ctx context.Context, t core.ReportType, path string) (string, error) {
	service, err := s.service(t)
	if err != nil {
		return "", err
	}
	return service.Find(ctx, path)
}

// Open coverage report with given path
func (s *Service) Open(ctx context.Context, t core.ReportType, path string) (io.Reader, error) {
	service, err := s.service(t)
	if err != nil {
		return nil, err
	}
	return service.Open(ctx, path)
}

// TrimFileNames for all files in coverage report
func (s *Service) TrimFileNames(ctx context.Context, report *core.CoverageReport, filters core.FileNameFilters) error {
	regexps := toRegexps(filters)
	for _, file := range report.Files {
		for _, regex := range regexps {
			file.Name = regex.ReplaceAllString(file.Name, "")
		}
		file.Name = strings.Trim(file.Name, "/")
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
