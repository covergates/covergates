package models

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/mock"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	// load sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func getDatabaseService(t *testing.T) (*gomock.Controller, core.DatabaseService) {
	ctrl := gomock.NewController(t)
	mockService := mock.NewMockDatabaseService(ctrl)
	mockService.EXPECT().Session().AnyTimes().Return(db.New())
	return ctrl, mockService
}

func TestMain(m *testing.M) {
	log.SetReportCaller(true)
	cwd, _ := os.Getwd()
	tempFile, err := ioutil.TempFile(cwd, "*.db")
	if err != nil {
		log.Fatal(err)
	}
	tempFile.Close()
	x, err := gorm.Open("sqlite3", tempFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	db = x
	if err := migrate(db); err != nil {
		log.Fatal(err)
	}
	exit := m.Run()
	defer os.Exit(exit)
	db.Close()
	os.Remove(tempFile.Name())
}
