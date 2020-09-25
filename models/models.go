package models

import (
	"github.com/covergates/covergates/core"
	"gorm.io/gorm"
)

var (
	tables []interface{}
)

type databaseService struct {
	db *gorm.DB
}

// NewDatabaseService with GORM
func NewDatabaseService(db *gorm.DB) core.DatabaseService {
	return &databaseService{
		db: db,
	}
}

func (store *databaseService) Session() *gorm.DB {
	return store.db.Session(&gorm.Session{})
}

func (store *databaseService) Migrate() error {
	return migrate(store.db)
}

func init() {
	tables = append(tables,
		&Report{},
		&ReportComment{},
		&Reference{},
		&Coverage{},
		&User{},
		&Repo{},
		&RepoSetting{},
		&RepoHook{},
		&OAuthToken{},
	)
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(tables...)
}
