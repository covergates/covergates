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
	provideReportStore,
	provideRepoStore,
)

func provideDatabaseService(db *gorm.DB) core.DatabaseService {
	return models.NewDatabaseService(db)
}

func provideUserStore(db core.DatabaseService) core.UserStore {
	return &models.UserStore{
		DB: db,
	}
}

func provideReportStore(db core.DatabaseService) core.ReportStore {
	return &models.ReportStore{
		DB: db,
	}
}

func provideRepoStore(db core.DatabaseService) core.RepoStore {
	return &models.RepoStore{
		DB: db,
	}
}
