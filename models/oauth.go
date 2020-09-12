package models

import (
	"fmt"
	"time"

	"github.com/covergates/covergates/core"
	"github.com/jinzhu/gorm"
)

// OAuthToken information
type OAuthToken struct {
	gorm.Model
	Name    string `gorm:"index"`
	Code    string `gorm:"index"`
	Access  string `gorm:"index"`
	Refresh string `gorm:"index"`
	Expires time.Time
	OwnerID uint
	Owner   *User `gorm:"foreignKey:OwnerID"`
	Data    []byte
}

// OAuthStore tokens
type OAuthStore struct {
	DB core.DatabaseService
}

// Create token
func (store *OAuthStore) Create(token *core.OAuthToken) error {
	if token.Owner == nil || token.Owner.Login == "" {
		return fmt.Errorf("require token owner")
	}
	session := store.DB.Session()
	user := &User{}
	if err := session.Where(&User{Login: token.Owner.Login}).First(user).Error; err != nil {
		return err
	}
	m := &OAuthToken{
		Name:    token.Name,
		Access:  token.Access,
		Code:    token.Code,
		Refresh: token.Refresh,
		Expires: token.Expires,
		Owner:   user,
		Data:    token.Data,
	}
	return session.Create(m).Error
}

// Find token with seed
func (store *OAuthStore) Find(token *core.OAuthToken) (*core.OAuthToken, error) {
	cond := &OAuthToken{
		Code:    token.Code,
		Access:  token.Access,
		Refresh: token.Refresh,
	}
	session := store.DB.Session()
	r := &OAuthToken{}
	if err := session.Preload("Owner").Where(cond).First(r).Error; err != nil {
		return nil, err
	}
	return r.toCoreOAuthToken(), nil
}

// List user's oauth tokens
func (store *OAuthStore) List(user *core.User) ([]*core.OAuthToken, error) {
	session := store.DB.Session()
	session.LogMode(true)
	u := &User{}
	if err := session.Where(&User{Login: user.Login}).First(u).Error; err != nil {
		return nil, err
	}
	var tokens []*OAuthToken
	if err := session.Preload("Owner").Where(&OAuthToken{OwnerID: u.ID}).Find(&tokens).Error; err != nil {
		return nil, err
	}
	result := make([]*core.OAuthToken, len(tokens))
	for i, token := range tokens {
		result[i] = token.toCoreOAuthToken()
	}
	return result, nil
}

// Delete token with seed
func (store *OAuthStore) Delete(token *core.OAuthToken) error {
	cond := &OAuthToken{
		Code:    token.Code,
		Access:  token.Access,
		Refresh: token.Refresh,
	}
	session := store.DB.Session()
	return session.Where(cond).Delete(OAuthToken{}).Error
}

func (token *OAuthToken) toCoreOAuthToken() *core.OAuthToken {
	return &core.OAuthToken{
		Name:    token.Name,
		Code:    token.Code,
		Access:  token.Access,
		Refresh: token.Refresh,
		Expires: token.Expires,
		Owner:   token.Owner.toCoreUser(),
		Data:    token.Data,
	}
}

// TableName for GORM
func (OAuthToken) TableName() string {
	return "oauth_token"
}
