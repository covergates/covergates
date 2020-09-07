package oauth

import (
	"context"
	"encoding/json"

	"github.com/covergates/covergates/core"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
)

type tokenStore struct {
	oauth2.TokenStore
	store core.OAuthStore
}

func (s *tokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	token := &core.OAuthToken{
		Data: data,
	}
	if code := info.GetCode(); code != "" {
		token.Code = code
		token.Expires = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn())
	} else {
		token.Access = info.GetAccess()
		token.Expires = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn())
		if refresh := info.GetRefresh(); refresh != "" {
			token.Refresh = refresh
			token.Expires = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
		}
	}
	return s.store.Create(token)
}

func (s *tokenStore) RemoveByCode(ctx context.Context, code string) error {
	return s.store.Delete(&core.OAuthToken{Code: code})
}

func (s *tokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return s.store.Delete(&core.OAuthToken{Access: access})
}

func (s *tokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return s.store.Delete(&core.OAuthToken{Refresh: refresh})
}

func (s *tokenStore) toTokenInfo(token *core.OAuthToken) oauth2.TokenInfo {
	var item models.Token
	if err := json.Unmarshal(token.Data, &item); err != nil {
		return nil
	}
	return &item
}

func (s *tokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	token, err := s.store.Find(&core.OAuthToken{Code: code})
	if err != nil {
		return nil, err
	}
	return s.toTokenInfo(token), nil
}

func (s *tokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	token, err := s.store.Find(&core.OAuthToken{Access: access})
	if err != nil {
		return nil, err
	}
	return s.toTokenInfo(token), nil
}

func (s *tokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	token, err := s.store.Find(&core.OAuthToken{Refresh: refresh})
	if err != nil {
		return nil, err
	}
	return s.toTokenInfo(token), nil
}
