package python

import (
	"context"
	"encoding/xml"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/service/common"
)

// CoverageService of python
type CoverageService struct{}

type coverage struct {
	XMLName  xml.Name `xml:"coverage"`
	Sources  []string `xml:"sources>source"`
	Packages []pkg    `xml:"packages>package"`
}

type pkg struct {
	Name    string  `xml:"name,attr"`
	Classes []class `xml:"classes>class"`
}

type class struct {
	Name     string  `xml:"filename,attr"`
	Lines    []line  `xml:"lines>line"`
	LineRate float64 `xml:"line-rate,attr"`
}

type line struct {
	Number int `xml:"number,attr"`
	Hits   int `xml:"hits,attr"`
}

// Report of python coverage
func (s *CoverageService) Report(ctx context.Context, reader io.Reader) (*core.CoverageReport, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	c := &coverage{}
	if err := xml.Unmarshal(data, c); err != nil {
		return nil, err
	}
	result := make([]*core.File, 0)
	for _, pkg := range c.Packages {
		files := pkg.toFiles(c.Source())
		result = append(result, files...)
	}
	report := &core.CoverageReport{
		Type:  core.ReportPython,
		Files: result,
	}
	report.StatementCoverage = report.ComputeStatementCoverage()
	return report, nil
}

//Find python coverage report in xml format
func (s *CoverageService) Find(ctx context.Context, path string) (string, error) {
	return common.FindReport(path, "coverage.xml")
}

// Open reader of python coverage report
func (s *CoverageService) Open(ctx context.Context, path string) (io.Reader, error) {
	return common.OpenFileReader(path)
}

func (c coverage) Source() string {
	return c.Sources[0]
}

func (p pkg) toFiles(parent string) []*core.File {
	files := make([]*core.File, len(p.Classes))
	for i, class := range p.Classes {
		statementHist := make([]*core.StatementHit, len(class.Lines))
		for j, line := range class.Lines {
			statementHist[j] = &core.StatementHit{
				Hits:       line.Hits,
				LineNumber: line.Number,
			}
		}
		files[i] = &core.File{
			Name:              filepath.Join(parent, class.Name),
			StatementHits:     statementHist,
			StatementCoverage: class.LineRate,
		}
	}
	return files
}
