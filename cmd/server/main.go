package main

import (
	"fmt"
	"os"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/routers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	// load sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func connectDatabase() *gorm.DB {
	x, _ := gorm.Open("sqlite3", "core.db")
	return x

}

// Run server
func Run(c *cli.Context) error {
	config, err := config.Environ()
	if err != nil {
		return err
	}
	db := connectDatabase()
	app, err := InitializeApplication(config, db)
	if err != nil {
		return err
	}
	r := gin.Default()
	app.routers.RegisterRoutes(r)
	app.db.Migrate()
	r.Run(fmt.Sprintf(":%s", config.Server.Port()))
	return nil
}

func main() {
	log.SetReportCaller(true)
	app := &cli.App{
		Name:   "codecover",
		Action: Run,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type application struct {
	routers *routers.Routers
	db      core.DatabaseService
}

func newApplication(
	routers *routers.Routers,
	db core.DatabaseService,
) application {
	return application{
		routers: routers,
		db:      db,
	}
}
