package models

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/mock"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func getDatabaseService(t *testing.T) (*gomock.Controller, core.DatabaseService) {
	ctrl := gomock.NewController(t)
	mockService := mock.NewMockDatabaseService(ctrl)
	mockService.EXPECT().Session().AnyTimes().Return(db.Session(&gorm.Session{}))
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
	x, err := gorm.Open(sqlite.Open(tempFile.Name()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	db = x
	if err := migrate(db); err != nil {
		log.Fatal(err)
	}
	exit := m.Run()
	defer os.Exit(exit)
	os.Remove(tempFile.Name())
}
