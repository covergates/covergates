package main

import (
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/models"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var storeSet = wire.NewSet(
	provideDatabaseService,
	provideUserStore,
)

func provideDatabaseService(db *gorm.DB) core.DatabaseService {
	return models.NewDatabaseService(db)
}

func provideUserStore(db core.DatabaseService) core.UserStore {
	return &models.UserStore{
		DB: db,
	}
}
