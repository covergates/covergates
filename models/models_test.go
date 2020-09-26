package models

import (
	"os"
	"testing"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/mock"
	"github.com/covergates/covergates/models/tests"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
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
	db = tests.OpenDatabase()
	if err := migrate(db); err != nil {
		log.Fatal(err)
	}
	exit := m.Run()
	defer os.Exit(exit)
	tests.CloseDatabase()
}
