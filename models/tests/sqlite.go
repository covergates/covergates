// +build !mysql,!postgres

package tests

import (
	"io/ioutil"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbName string
)

// OpenDatabase with sqlite
func OpenDatabase() *gorm.DB {
	cwd, _ := os.Getwd()
	tempFile, err := ioutil.TempFile(cwd, "*.db")
	if err != nil {
		log.Fatal(err)
	}
	tempFile.Close()
	dbName = tempFile.Name()
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// CloseDatabase remove temporary files
func CloseDatabase() {
	os.Remove(dbName)
}
