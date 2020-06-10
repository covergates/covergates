//+build wireinject

package main

import (
	"github.com/code-devel-cover/CodeCover/config"
	"github.com/google/wire"
)

func InitializeApplication(config *config.Config) (application, error) {
	wire.Build(
		routerSet,
		newApplication,
	)
	return application{}, nil
}
