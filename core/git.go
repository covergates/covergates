package core

import "context"

//go:generate mockgen -package mock -destination ../mock/git_mock.go . GitRepository,Git

// GitRepository which is cloned from SCM
type GitRepository interface {
	ListAllFiles(commit string) ([]string, error)
}

// Git interact with SCM with plain git commands
type Git interface {
	Clone(ctx context.Context, URL, token string) (GitRepository, error)
}
