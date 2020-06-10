package core

import "github.com/drone/go-scm/scm"

type SCMClientService interface {
	Client(scm SCMProvider) *scm.Client
}
