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

	// load drivers
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Version of covergates server
var Version = "0.0"

func connectDatabase(cfg *config.Config) *gorm.DB {
	var x *gorm.DB
	var err error
	switch cfg.Database.Driver {
	case "sqlite3":
		x, err = gorm.Open(cfg.Database.Driver, cfg.Database.Name)
	case "postgres":
		x, err = gorm.Open(
			cfg.Database.Driver,
			fmt.Sprintf(
				"host=%s port=%s user=%s password=%s database=%s",
				cfg.Database.Host,
				cfg.Database.Port,
				cfg.Database.User,
				cfg.Database.Password,
				cfg.Database.Name,
			))
	case "cloudrun":
		x, err = gorm.Open(
			"postgres",
			fmt.Sprintf(
				"user=%s password=%s database=%s host=%s/%s",
				cfg.CloudRun.User,
				cfg.CloudRun.Password,
				cfg.CloudRun.Name,
				cfg.CloudRun.Socket,
				cfg.CloudRun.Instance))
	default:
		log.Fatal("database driver not support")
	}
	if err != nil {
		log.Fatal(err)
	}
	return x
}

// Run server
func Run(c *cli.Context) error {
	config, err := config.Environ()
	if err != nil {
		return err
	}
	db := connectDatabase(config)
	app, err := InitializeApplication(config, db)
	if err != nil {
		return err
	}
	if config.Database.AutoMigrate {
		go func() {
			app.db.Migrate()
			log.Println("migration done")
		}()
	}
	r := gin.Default()
	app.routers.RegisterRoutes(r)
	r.Run(fmt.Sprintf(":%s", config.Server.Port()))
	return nil
}

func main() {
	log.SetReportCaller(true)
	app := &cli.App{
		Name:    "codecover",
		Version: Version,
		Action:  Run,
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
