package models

import "github.com/jinzhu/gorm"

var (
	db        *gorm.DB
	tables    []interface{}
	connected bool
)

func init() {

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
func ConnectDB(db *gorm.DB) error {
	if connected {
		return nil
	}
	if err := migrate(); err != nil {
		return err
	}
	connected = true
	return nil
}
