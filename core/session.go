package core

import "github.com/gin-gonic/gin"

// Session of user login for Gin
type Session interface {
	CreateUser(c *gin.Context, user *User) error
	GetUser(c *gin.Context) *User
	StartBindUser(c *gin.Context) error
	EndBindUser(c *gin.Context) error
	ShouldBindUser(c *gin.Context) bool
	Clear(c *gin.Context) error
}
