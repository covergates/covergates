package main

import (
	"log"
	"os"

	"github.com/covergates/covergates/cmd/cli/comment"
	"github.com/covergates/covergates/cmd/cli/upload"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "covergate",
		Commands: []*cli.Command{
			upload.Command,
			comment.Command,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "url",
				Value:   "http://localhost:5900/api/v1",
				Usage:   "api service url",
				EnvVars: []string{"API_URL"},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
