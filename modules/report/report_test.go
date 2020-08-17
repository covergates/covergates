package report

import (
	"io/ioutil"
	"testing"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

const expectMarkdown = `### [Coverage: 50.0%](http://localhost/report/github/space/name?ref=commit)

||File|Coverage|
|--|--|--------|
|:arrow_up_small:|A|1.00|
|:arrow_down_small:|B|0.00|
||C|0.50|
`

const expectMarkdownNoTarget = `### [Coverage: 50.0%](http://localhost/report/github/space/name?ref=commit)

||File|Coverage|
|--|--|--------|
|:arrow_up_small:|A|1.00|
||B|0.00|
|:arrow_up_small:|C|0.50|
`

func TestMarkdownReport(t *testing.T) {
	source := &core.Report{
		ReportID: "report_id",
		Commit:   "commit",
		Coverages: []*core.CoverageReport{
			{
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
		},
	}

	target := &core.Report{
		Coverages: []*core.CoverageReport{
			{
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
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepoStore(ctrl)
	mockRepo.EXPECT().Find(gomock.Eq(
		&core.Repo{ReportID: source.ReportID},
	)).AnyTimes().Return(&core.Repo{
		Name:      "name",
		NameSpace: "space",
		SCM:       core.Github,
	}, nil)

	service := &Service{
		Config: &config.Config{
			Server: config.Server{
				Addr: "http://localhost",
			},
		},
		RepoStore: mockRepo,
	}

	reader, err := service.MarkdownReport(source, target)
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(expectMarkdown, string(data)); diff != "" {
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
	if diff := cmp.Diff(expectMarkdownNoTarget, string(data)); diff != "" {
		t.Log(diff)
		t.Fail()
	}
}
