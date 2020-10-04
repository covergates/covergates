package clover

import (
	"context"
	"encoding/xml"
	"io"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/service/common"
)

const (
	typeStmt   = "stmt"
	typeMethod = "method"
)

// CoverageService for clover report
type CoverageService struct {
}

type coverage struct {
	XMLName  xml.Name  `xml:"coverage"`
	Packages []pkg     `xml:"project>package"`
	Files    fileSlice `xml:"project>file"`
}

type pkg struct {
	Name  string    `xml:"name,attr"`
	Files fileSlice `xml:"file"`
}

type fileSlice []file

type file struct {
	Name   string `xml:"name,attr"`
	Lines  []line `xml:"line"`
	Metric metric `xml:"metrics"`
}

type metric struct {
	Statements        int `xml:"statements,attr"`
	CoveredStatements int `xml:"coveredstatements,attr"`
}

type line struct {
	Num   int    `xml:"num,attr"`
	Type  string `xml:"type,attr"`
	Count int    `xml:"count,attr"`
}

// Report for clover
func (s *CoverageService) Report(ctx context.Context, data io.Reader) (*core.CoverageReport, error) {
	cov := &coverage{}
	if err := xml.NewDecoder(data).Decode(cov); err != nil {
		return nil, err
	}
	return cov.toCoverageReport()
}

// Find clover report, if path is file, return it directly
func (s *CoverageService) Find(ctx context.Context, path string) (string, error) {
	return common.FindReport(path, "clover.xml")
}

// Open reader of clover coverage report
func (s *CoverageService) Open(ctx context.Context, path string) (io.Reader, error) {
	return common.OpenFileReader(path)
}

func (c *coverage) toCoverageReport() (*core.CoverageReport, error) {
	files := make([]*core.File, 0)
	for _, pkg := range c.Packages {
		s, err := pkg.Files.toCoreFiles()
		if err != nil {
			return nil, err
		}
		files = append(files, s...)
	}
	s, err := c.Files.toCoreFiles()
	if err != nil {
		return nil, err
	}
	files = append(files, s...)
	report := &core.CoverageReport{
		Files: files,
		Type:  core.ReportClover,
	}
	report.StatementCoverage = report.ComputeStatementCoverage()
	return report, nil
}

func (f *file) toFile() (*core.File, error) {
	hits := make([]*core.StatementHit, 0)
	for _, line := range f.Lines {
		if line.Type != typeStmt {
			continue
		}
		hits = append(hits, &core.StatementHit{
			Hits:       (line.Count),
			LineNumber: line.Num,
		})
	}
	coverage := 0.0
	if f.Metric.Statements > 0 {
		coverage = float64(f.Metric.CoveredStatements) / float64(f.Metric.Statements)

	}
	return &core.File{
		Name:              f.Name,
		StatementHits:     hits,
		StatementCoverage: coverage,
	}, nil
}

func (s fileSlice) toCoreFiles() ([]*core.File, error) {
	var err error
	files := make([]*core.File, len(s))
	for i, file := range s {
		if files[i], err = file.toFile(); err != nil {
			return nil, err
		}
	}
	return files, nil
}
