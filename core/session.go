package core

import "github.com/gin-gonic/gin"

// Session of user login for Gin
type Session interface {
	Create(c *gin.Context, user *User) error
	Get(c *gin.Context) *User
	Clear(c *gin.Context) error
}
