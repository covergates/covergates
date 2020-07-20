package models

import (
	"fmt"

	"github.com/code-devel-cover/CodeCover/core"
)

type errNotSupportedSCM struct {
	scm core.SCMProvider
}

func (e *errNotSupportedSCM) Error() string {
	return fmt.Sprintf("%s is not supported", string(e.scm))
}

// IsErrNotSupportedSCM check
func IsErrNotSupportedSCM(err error) bool {
	_, ok := err.(*errNotSupportedSCM)
	return ok
}
