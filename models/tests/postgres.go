// +build postgres

package tests

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OpenDatabase with postgres
func OpenDatabase() *gorm.DB {
	db, err := gorm.Open(
		postgres.Open("host=127.0.0.1 port=5432 user=covergates password=covergates database=covergates"),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// CloseDatabase does nothing
func CloseDatabase() {
	return
}
