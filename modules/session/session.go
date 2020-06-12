package session

import (
	"encoding/gob"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	gob.Register(&core.User{})
}

// Session stores records for Gin
type Session struct {
}

const (
	keyUser = "user"
)

// Create a user record
func (s *Session) Create(c *gin.Context, user *core.User) error {
	session := sessions.Default(c)
	session.Set(keyUser, user)
	return session.Save()
}

// Get a user
func (s *Session) Get(c *gin.Context) *core.User {
	session := sessions.Default(c)
	data := session.Get(keyUser)
	user, ok := data.(*core.User)
	if !ok {
		return &core.User{}
	}
	return user
}

// Clear all session records
func (s *Session) Clear(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	return session.Save()
}
