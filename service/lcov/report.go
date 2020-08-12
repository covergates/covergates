package lcov

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/service/common"
)

// Reference: https://github.com/linux-test-project/lcov/blob/master/bin/geninfo

// CoverageService of lcov
type CoverageService struct{}

const (
	labelEndOfRecord = "end_of_record"
	labelSF          = "SF"
	labelDA          = "DA"
)

var (
	errFormat = errors.New("invalid lcov format")
)

// Report of lcov
func (s *CoverageService) Report(ctx context.Context, reader io.Reader) (*core.CoverageReport, error) {

	scanner := bufio.NewScanner(reader)
	files := make([]*core.File, 0)
	var cf *core.File
	for scanner.Scan() {
		line := scanner.Text()
		if line == labelEndOfRecord {
			if cf == nil {
				return nil, errFormat
			}
			cf.StatementCoverage = common.ComputeStatementCoverage(cf.StatementHits)
			files = append(files, cf)
			cf = nil
			continue
		} else if line == "" {
			continue
		}
		tokens := strings.SplitN(line, ":", 2)
		if len(tokens) != 2 {
			return nil, errFormat
		}
		if tokens[0] == labelSF {
			cf = &core.File{
				Name:          tokens[1],
				StatementHits: make([]*core.StatementHit, 0),
			}
		} else if tokens[0] == labelDA {
			tokens = strings.SplitN(tokens[1], ",", 2)
			if len(tokens) != 2 {
				return nil, errFormat
			}
			ln, err := strconv.Atoi(tokens[0])
			if err != nil {
				return nil, err
			}
			hit, err := strconv.Atoi(tokens[1])
			if err != nil {
				return nil, err
			}
			cf.StatementHits = append(cf.StatementHits, &core.StatementHit{
				LineNumber: ln,
				Hits:       hit,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	report := &core.CoverageReport{
		Files: files,
		Type:  core.ReportLCOV,
	}
	report.StatementCoverage = report.ComputeStatementCoverage()
	return report, nil
}

//Find lcov coverage report
func (s *CoverageService) Find(ctx context.Context, path string) (string, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if fi.Mode().IsRegular() {
		return path, nil
	}
	return common.FindReport(path, "lcov.info")
}

// Open reader of lcov coverage report
func (s *CoverageService) Open(ctx context.Context, path string) (io.Reader, error) {
	return common.OpenFileReader(path)
}
