package models

import (
	"github.com/jinzhu/gorm"
	// load sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db        *gorm.DB
	tables    []interface{}
	connected bool
)

func init() {
	tables = append(tables, &Report{})
}

func migrate() error {
	for _, table := range tables {
		if err := db.AutoMigrate(table).Error; err != nil {
			return err
		}
	}
	return nil
}

// ConnectDB with gorm
func ConnectDB(x *gorm.DB) error {
	if connected {
		return nil
	}
	db = x
	if err := migrate(); err != nil {
		return err
	}
	connected = true
	return nil
}
