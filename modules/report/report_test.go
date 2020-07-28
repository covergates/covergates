package report

import (
	"io/ioutil"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/google/go-cmp/cmp"
)

const expectMarkdown = `### Coverage: 80.0%

||File|Coverage|
|--|--|--------|
|:arrow_up_small:|A|1.00|
|:arrow_down_small:|B|0.00|
||C|0.50|
`

const expectMarkdownNoTarget = `### Coverage: 80.0%

||File|Coverage|
|--|--|--------|
|:arrow_up_small:|A|1.00|
||B|0.00|
|:arrow_up_small:|C|0.50|
`

func TestMarkdownReport(t *testing.T) {
	source := &core.Report{
		Coverage: &core.CoverageReport{
			StatementCoverage: 0.8,
			Files: []*core.File{
				{
					Name:              "A",
					StatementCoverage: 1.0,
				},
				{
					Name:              "B",
					StatementCoverage: 0,
				},
				{
					Name:              "C",
					StatementCoverage: 0.5,
				},
			},
		},
	}
	target := &core.Report{
		Coverage: &core.CoverageReport{
			StatementCoverage: 0.8,
			Files: []*core.File{
				{
					Name:              "A",
					StatementCoverage: 0.8,
				},
				{
					Name:              "B",
					StatementCoverage: 0.8,
				},
				{
					Name:              "C",
					StatementCoverage: 0.5,
				},
			},
		},
	}
	service := &Service{}
	reader, err := service.MarkdownReport(source, target)
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(data), expectMarkdown); diff != "" {
		t.Log(diff)
		t.Fail()
	}

	reader, err = service.MarkdownReport(source, &core.Report{})
	if err != nil {
		t.Fatal(err)
	}
	data, err = ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(data), expectMarkdownNoTarget); diff != "" {
		t.Log(diff)
		t.Fail()
	}
}
