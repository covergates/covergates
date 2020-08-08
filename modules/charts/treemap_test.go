package charts

import (
	"os"
	"testing"

	"github.com/covergates/covergates/core"
)

func TestCoverDiffTreeMap(t *testing.T) {
	o := &core.Report{
		Coverages: []*core.CoverageReport{
			{
				Files: []*core.File{
					{
						Name:              "A",
						StatementCoverage: 0.5,
						StatementHits:     make([]*core.StatementHit, 20),
					},
					{
						Name:              "B",
						StatementCoverage: 0.8,
						StatementHits:     make([]*core.StatementHit, 10),
					},
					{
						Name:              "C",
						StatementCoverage: 1.0,
						StatementHits:     make([]*core.StatementHit, 50),
					},
				},
			},
		},
	}

	n := &core.Report{
		Coverages: []*core.CoverageReport{
			{
				Files: []*core.File{
					{
						Name:              "A",
						StatementCoverage: 0.6,
						StatementHits:     make([]*core.StatementHit, 20),
					},
					{
						Name:              "B",
						StatementCoverage: 0.4,
						StatementHits:     make([]*core.StatementHit, 10),
					},
					{
						Name:              "C",
						StatementCoverage: 1.0,
						StatementHits:     make([]*core.StatementHit, 50),
					},
				},
			},
		},
	}
	m := NewCoverageDiffTreeMap(o, n)
	file, err := os.Create("treemap.svg")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()
	if err := m.Render(file); err != nil {
		t.Error(err)
		return
	}
}
