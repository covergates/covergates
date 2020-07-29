package core

import "context"

// HookService manages and resolves webhooks
type HookService interface {
	Create(ctx context.Context, repo *Repo) error
	Delete(ctx context.Context, repo *Repo) error
	Resolve(ctx context.Context, repo *Repo, hook HookEvent) error
}
