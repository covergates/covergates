package upload

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/modules/util"
	"github.com/code-devel-cover/CodeCover/service/coverage"
	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
)

// Command for upload report
var Command = &cli.Command{
	Name:      "upload",
	Usage:     "upload coverage report",
	ArgsUsage: "report",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "report",
			Usage:    "report id",
			EnvVars:  []string{"REPORT_ID"},
			Value:    "",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "type",
			Usage:    "report type",
			Value:    "",
			Required: true,
		},
	},
	Action: upload,
}

type head struct {
	commit string
	branch string
}

func upload(c *cli.Context) error {
	if c.NArg() <= 0 {
		cli.ShowCommandHelp(c, "upload")
		return fmt.Errorf("report path is required")
	}

	data, err := findReportData(c.Context, c.String("type"), c.Args().First())
	if err != nil {
		return err
	}

	h, err := repoHead()
	if err != nil {
		return err
	}

	form := util.FormData{
		"type":   c.String("type"),
		"commit": h.commit,
		"branch": h.branch,
		"file": util.FormFile{
			Name: "report",
			Data: data,
		},
	}

	url := fmt.Sprintf(
		"%s/reports/%s",
		c.String("url"),
		c.String("report"),
	)

	respond, err := util.PostForm(url, form)
	if err != nil {
		return err
	}

	defer respond.Body.Close()
	text, err := ioutil.ReadAll(respond.Body)
	if respond.StatusCode >= 400 {
		log.Fatal(string(text))
	}
	log.Println(string(text))
	return nil
}

func findReportData(ctx context.Context, reportType, path string) ([]byte, error) {
	t := core.ReportType(reportType)
	service := &coverage.Service{}
	report, err := service.Find(ctx, t, path)
	if err != nil {
		return nil, err
	}
	r, err := service.Open(ctx, t, report)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

func repoHead() (*head, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpenWithOptions(cwd, &git.PlainOpenOptions{
		DetectDotGit: true,
	})

	if err != nil {
		return nil, err
	}

	h, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit := h.Hash().String()
	branch := ""
	if h.Name().IsBranch() {
		branch = h.Name().Short()
	}
	return &head{
		branch: branch,
		commit: commit,
	}, nil
}
