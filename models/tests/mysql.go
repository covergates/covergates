// +build mysql

package tests

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OpenDatabase with mysql
func OpenDatabase() *gorm.DB {
	db, err := gorm.Open(
		mysql.Open("covergates:covergates@tcp(127.0.0.1:3306)/covergates?parseTime=true&loc=Local"),
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
