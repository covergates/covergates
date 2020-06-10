package main

import (
	"log"
	"os"

	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/routers"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name: "codecover",
		Action: func(c *cli.Context) error {
			config := &config.Config{
				Github: config.Github{
					ClientID:     "a150e893154bafde8a00",
					ClientSecret: "59a3f97b6e7569d0b6898bc5fb2e84f93e64113d",
					Server:       "https://github.com",
				},
			}
			app, err := InitializeApplication(config)
			if err != nil {
				return err
			}
			r := gin.Default()
			app.routers.RegisterRoutes(r)
			r.Run(":5900")
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type application struct {
	routers *routers.Routers
}

func newApplication(
	routers *routers.Routers,
) application {
	return application{
		routers: routers,
	}
}
