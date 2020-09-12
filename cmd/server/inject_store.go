package main

import (
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/models"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var storeSet = wire.NewSet(
	provideDatabaseService,
	provideUserStore,
	provideReportStore,
	provideRepoStore,
	provideOAuthStore,
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

func provideOAuthStore(db core.DatabaseService) core.OAuthStore {
	return &models.OAuthStore{
		DB: db,
	}
}
