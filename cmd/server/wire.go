//+build wireinject

package main

import (
	"github.com/code-devel-cover/CodeCover/config"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

func InitializeApplication(config *config.Config, db *gorm.DB) (application, error) {
	wire.Build(
		serviceSet,
		storeSet,
		routerSet,
		newApplication,
	)
	return application{}, nil
}
