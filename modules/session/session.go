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
	keyUser    = "user"
	keyBinding = "bindUser"
)

// CreateUser in session record
func (s *Session) CreateUser(c *gin.Context, user *core.User) error {
	session := sessions.Default(c)
	session.Set(keyUser, user)
	return session.Save()
}

// GetUser from session
func (s *Session) GetUser(c *gin.Context) *core.User {
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

// StartBindUser remarks a user is going to bind new SCM account in session
func (s *Session) StartBindUser(c *gin.Context) error {
	session := sessions.Default(c)
	session.Set(keyBinding, true)
	return session.Save()
}

// EndBindUser remarks server finishes bind new SCM account in session
func (s *Session) EndBindUser(c *gin.Context) error {
	session := sessions.Default(c)
	session.Delete(keyBinding)
	return session.Save()
}

// ShouldBindUser tells if the new SCM login should be with a user
func (s *Session) ShouldBindUser(c *gin.Context) bool {
	session := sessions.Default(c)
	shouldBind, ok := session.Get(keyBinding).(bool)
	if !ok || !shouldBind {
		return false
	}
	return true
}
