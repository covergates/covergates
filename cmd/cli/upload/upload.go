package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/modules/git"
	"github.com/covergates/covergates/modules/util"
	"github.com/covergates/covergates/service/coverage"
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
		&cli.StringFlag{
			Name:     "branch",
			Usage:    "branch to upload the report",
			EnvVars:  []string{"GITHUB_HEAD_REF", "DRONE_SOURCE_BRANCH"},
			Value:    "",
			Required: false,
		},
	},
	Action: upload,
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

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	git := &git.Service{}
	repo, err := git.PlainOpen(c.Context, cwd)
	if err != nil {
		return err
	}

	branch := c.String("branch")
	if branch == "" {
		branch = repo.Branch()
	}

	files, err := repo.ListAllFiles(repo.HeadCommit())
	if err != nil {
		return err
	}
	filesData, err := json.Marshal(files)
	if err != nil {
		return err
	}

	form := util.FormData{
		"type":   c.String("type"),
		"commit": repo.HeadCommit(),
		"ref":    branch,
		"files":  string(filesData),
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

	log.Printf("upload commit %s, %s\n", repo.HeadCommit(), c.String("type"))

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
