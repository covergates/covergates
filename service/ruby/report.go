package ruby

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/service/common"
)

type resultType string

type resultSet map[resultType]*result

type result struct {
	Coverage  map[string]*coverage `json:"coverage"`
	Timestamp int                  `json:"timestamp"`
}

type coverage struct {
	Lines []interface{} `json:"lines"`
}

// CoverageService of ruby
type CoverageService struct{}

const (
	rspec resultType = "RSpec"
)

var (
	errResultTypeNotSupport = fmt.Errorf("report type not support (only support RSpec)")
)

// Report of ruby (rspec) coverage
func (s *CoverageService) Report(ctx context.Context, reader io.Reader) (*core.CoverageReport, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var results resultSet
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	rspecResult, ok := results[rspec]
	if !ok {
		return nil, errResultTypeNotSupport
	}

	files := make([]*core.File, 0)
	for name, cov := range rspecResult.Coverage {
		hits := make([]*core.StatementHit, 0)
		for i, hit := range cov.Lines {
			if hit == nil {
				continue
			}
			switch v := hit.(type) {
			case int:
				hits = append(hits, &core.StatementHit{
					Hits:       v,
					LineNumber: i + 1,
				})
			case float64:
				hits = append(hits, &core.StatementHit{
					Hits:       int(v),
					LineNumber: i + 1,
				})
			}
		}
		file := &core.File{
			Name:              name,
			StatementHits:     hits,
			StatementCoverage: common.ComputeStatementCoverage(hits),
		}
		files = append(files, file)
	}

	coverageReport := &core.CoverageReport{
		Files: files,
		Type:  core.ReportRuby,
	}
	coverageReport.StatementCoverage = coverageReport.ComputeStatementCoverage()
	return coverageReport, nil
}

//Find ruby coverage report in JSON format
func (s *CoverageService) Find(ctx context.Context, path string) (string, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if fi.Mode().IsRegular() {
		return path, nil
	}
	return common.FindReport(path, ".resultset.json")
}

// Open reader of ruby coverage report
func (s *CoverageService) Open(ctx context.Context, path string) (io.Reader, error) {
	return common.OpenFileReader(path)
}
