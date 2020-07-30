package request

import (
	"github.com/covergates/covergates/core"
	"github.com/gin-gonic/gin"
)

const (
	keyUser = "user"
)

// WithUser context
func WithUser(c *gin.Context, user *core.User) {
	c.Set(keyUser, user)
}

// UserFrom context
func UserFrom(c *gin.Context) (*core.User, bool) {
	data, ok := c.Get(keyUser)
	if !ok {
		return nil, false
	}
	user, ok := data.(*core.User)
	if !ok {
		return nil, false
	}
	return user, true
}
