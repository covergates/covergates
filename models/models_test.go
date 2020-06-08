package models

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestMain(m *testing.M) {
	cwd, _ := os.Getwd()
	tempFile, err := ioutil.TempFile(cwd, "*.db")
	if err != nil {
		log.Fatal(err)
	}
	tempFile.Close()
	db, err := gorm.Open("sqlite3", tempFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	if err := ConnectDB(db); err != nil {
		log.Fatal(err)
	}
	exit := m.Run()
	defer os.Exit(exit)
	db.Close()
	os.Remove(tempFile.Name())
}
