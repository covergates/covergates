package core

import "github.com/jinzhu/gorm"

//go:generate mockgen -package mock -destination ../mock/store_mock.go . DatabaseService

type DatabaseService interface {
	Session() *gorm.DB
	Migrate() error
}
