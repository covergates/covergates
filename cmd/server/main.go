package main

import (
	"fmt"
	"os"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/routers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Version of covergates server
var Version = "0.0"

func connectDatabase(cfg *config.Config) *gorm.DB {
	var x *gorm.DB
	var err error
	switch cfg.Database.Driver {
	case "sqlite3":
		x, err = gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
	case "postgres":
		x, err = gorm.Open(
			postgres.Open(
				fmt.Sprintf(
					"host=%s port=%s user=%s password=%s database=%s",
					cfg.Database.Host,
					cfg.Database.Port,
					cfg.Database.User,
					cfg.Database.Password,
					cfg.Database.Name,
				)), &gorm.Config{})
	case "mysql":
		x, err = gorm.Open(
			mysql.Open(
				fmt.Sprintf(
					"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
					cfg.Database.User,
					cfg.Database.Password,
					cfg.Database.Host,
					cfg.Database.Port,
					cfg.Database.Name,
				)), &gorm.Config{},
		)
	case "cloudrun":
		x, err = gorm.Open(
			postgres.Open(
				fmt.Sprintf(
					"user=%s password=%s database=%s host=%s/%s",
					cfg.CloudRun.User,
					cfg.CloudRun.Password,
					cfg.CloudRun.Name,
					cfg.CloudRun.Socket,
					cfg.CloudRun.Instance,
				)), &gorm.Config{})
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
