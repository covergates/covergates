package models

import (
	"github.com/jinzhu/gorm"
)

var (
	tables []interface{}
)

type databaseService struct {
	db *gorm.DB
}

func NewDatabaseService(db *gorm.DB) *databaseService {
	return &databaseService{
		db: db,
	}
}

func (store *databaseService) Session() *gorm.DB {
	return store.db.New()
}

func (store *databaseService) Migrate() error {
	return migrate(store.db)
}

func init() {
	tables = append(tables,
		&Report{},
		&User{},
		&Repo{},
	)
}

func migrate(db *gorm.DB) error {
	for _, table := range tables {
		if err := db.AutoMigrate(table).Error; err != nil {
			return err
		}
	}
	return nil
}
