package oauth

import (
	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

// Service of OAuth
type Service struct {
	server *server.Server
}

// NewService for OAuth
func NewService(config *config.Config, store core.OAuthStore) *Service {
	return &Service{
		server: newOAuthServer(config, store),
	}
}

func newOAuthServer(config *config.Config, oauthStore core.OAuthStore) *server.Server {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(
		&tokenStore{
			store: oauthStore,
		},
		nil,
	)

	clientStore := store.NewClientStore()
	clientStore.Set(
		config.Server.OAuthClient,
		&models.Client{
			ID:     config.Server.OAuthClient,
			Secret: config.Server.Secret,
			Domain: config.Server.Addr,
		},
	)

	srv := server.NewDefaultServer(manager)

	return srv
}
