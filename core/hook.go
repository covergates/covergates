package core

import "context"

//go:generate mockgen -package mock -destination ../mock/hook_mock.go . HookService

// HookService manages and resolves webhooks
type HookService interface {
	Create(ctx context.Context, repo *Repo) error
	Delete(ctx context.Context, repo *Repo) error
	Resolve(ctx context.Context, repo *Repo, hook HookEvent) error
}
