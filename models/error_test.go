package models

import (
	"fmt"
	"testing"

	"github.com/covergates/covergates/core"
)

func TestErrNotSupportSCM(t *testing.T) {
	var e error
	e = &errNotSupportedSCM{scm: core.Gitea}
	if !IsErrNotSupportedSCM(e) {
		t.Fail()
	}
	e = fmt.Errorf("Fake error")
	if IsErrNotSupportedSCM(e) {
		t.Fail()
	}
}
