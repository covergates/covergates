package oauth

import (
	"context"
	"net/http"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	log "github.com/sirupsen/logrus"
)

type key string

const (
	userKey key = "OAuthUser"
	nameKey key = "TokenName"
)

// ErrTokenOwnerNotFound in context
var ErrTokenOwnerNotFound = errors.New("requires token owner")

// Service of OAuth
type Service struct {
	config     *config.Config
	server     *server.Server
	oauthStore core.OAuthStore
	userStore  core.UserStore
}

// NewService for OAuth
func NewService(
	config *config.Config,
	oauthStore core.OAuthStore,
	userStore core.UserStore,
) *Service {
	return &Service{
		config:     config,
		server:     newOAuthServer(config, oauthStore),
		oauthStore: oauthStore,
		userStore:  userStore,
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

	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)

	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Errorln(err)
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Errorln(re.Error)
		return
	})

	return srv
}

// CreateToken with name. Required context with User.
func (s *Service) CreateToken(ctx context.Context, name string) (*core.OAuthToken, error) {
	ctx = withTokenName(ctx, name)
	token, err := s.server.Manager.GenerateAccessToken(
		ctx,
		oauth2.ClientCredentials,
		&oauth2.TokenGenerateRequest{
			ClientID:     s.config.Server.OAuthClient,
			ClientSecret: s.config.Server.Secret,
		},
	)
	if err != nil {
		return nil, err
	}
	return s.oauthStore.Find(&core.OAuthToken{Code: token.GetCode()})
}

// DeleteToken with access code
func (s *Service) DeleteToken(ctx context.Context, token *core.OAuthToken) error {
	return s.server.Manager.RemoveAccessToken(ctx, token.Access)
}

// ListTokens of token owner
func (s *Service) ListTokens(ctx context.Context) ([]*core.OAuthToken, error) {
	user, ok := getUser(ctx)
	if !ok {
		return nil, ErrTokenOwnerNotFound
	}
	return s.oauthStore.List(user)
}

// Validate Bearer token and return token's User if exist
func (s *Service) Validate(r *http.Request) (*core.User, error) {
	token, err := s.server.ValidationBearerToken(r)
	if err != nil {
		return nil, err
	}
	userID := token.GetUserID()
	return s.userStore.FindByLogin(userID)
}

// WithUser context
func (s *Service) WithUser(ctx context.Context, user *core.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	user, ok := getUser(r.Context())
	if !ok {
		return "", nil
	}
	return user.Login, nil
}

func getUser(ctx context.Context) (*core.User, bool) {
	u, ok := ctx.Value(userKey).(*core.User)
	return u, ok
}

func withTokenName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, nameKey, name)
}

func getTokenName(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(nameKey).(string)
	return name, ok
}
