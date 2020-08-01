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

func main() {
	app := &cli.App{
		Name:    "covergate",
		Version: Version,
		Commands: []*cli.Command{
			upload.Command,
			comment.Command,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "url",
				Value:   CoverGatesAPI,
				Usage:   "api service url",
				EnvVars: []string{"API_URL"},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
