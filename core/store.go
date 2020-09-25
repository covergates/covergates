package core

import "gorm.io/gorm"

//go:generate mockgen -package mock -destination ../mock/store_mock.go . DatabaseService

// DatabaseService provides database operations with GORM
type DatabaseService interface {
	Session() *gorm.DB
	Migrate() error
}
