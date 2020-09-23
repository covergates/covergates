package main

import (
	"log"
	"os"

	"github.com/covergates/covergates/cmd/cli/comment"
	"github.com/covergates/covergates/cmd/cli/upload"
	"github.com/urfave/cli/v2"
)

var (
	// CoverGatesAPI to covergates API URL
	CoverGatesAPI = "http://localhost:8080/api/v1"
	// Version of cli
	Version = "0.0"
)

var app = &cli.App{
	Name:    "covergate",
	Version: Version,
	Commands: []*cli.Command{
		upload.Command,
		comment.Command,
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "token",
			Usage:   "provide OAuth token for API",
			EnvVars: []string{"GATES_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "url",
			Value:   CoverGatesAPI,
			Usage:   "api service url",
			EnvVars: []string{"API_URL"},
		},
	},
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
