package core

import (
	"github.com/drone/go-login/login"
)

// LoginMiddleware for Gin
type LoginMiddleware interface {
	Handler(scm SCMProvider) login.Middleware
}
