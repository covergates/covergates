package core

import (
	"github.com/drone/go-login/login"
)

type LoginMiddleware interface {
	Handler(scm SCMProvider) login.Middleware
}
