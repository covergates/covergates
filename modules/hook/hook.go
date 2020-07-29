package hook

import (
	"context"
	"errors"

	"github.com/code-devel-cover/CodeCover/core"
)

var errHookEventNotSupport = errors.New("hook event not support")

// Service of webhook resolve and management
type Service struct {
	SCM           core.SCMService
	RepoStore     core.RepoStore
	ReportStore   core.ReportStore
	ReportService core.ReportService
}

// Create a webhook to repository. If existed webhook found, it will be removed first
func (s *Service) Create(ctx context.Context, repo *core.Repo) error {
	user, err := s.RepoStore.Creator(repo)
	if err != nil {
		return err
	}
	client, err := s.SCM.Client(repo.SCM)
	if err != nil {
		return err
	}
	if hook, err := s.RepoStore.FindHook(repo); hook != nil && err == nil {
		client.Repositories().RemoveHook(ctx, user, repo.FullName(), hook)
	}
	hook, err := client.Repositories().CreateHook(ctx, user, repo.FullName())
	if err != nil {
		return err
	}
	return s.RepoStore.UpdateHook(repo, hook)
}

// Delete a repository webhook
func (s *Service) Delete(ctx context.Context, repo *core.Repo) error {
	user, err := s.RepoStore.Creator(repo)
	if err != nil {
		return err
	}
	client, err := s.SCM.Client(repo.SCM)
	if err != nil {
		return err
	}
	hook, err := s.RepoStore.FindHook(repo)
	if err != nil {
		return err
	}
	return client.Repositories().RemoveHook(ctx, user, repo.FullName(), hook)
}

// Resolve webhook event from the SCM
func (s *Service) Resolve(ctx context.Context, repo *core.Repo, hook core.HookEvent) error {
	if hook == nil {
		return errHookEventNotSupport
	}
	switch event := hook.(type) {
	case *core.PullRequestHook:
		return s.resolvePullRequest(ctx, repo, event)
	}
	return nil
}

func (s *Service) resolvePullRequest(ctx context.Context, repo *core.Repo, hook *core.PullRequestHook) error {
	setting, err := s.RepoStore.Setting(repo)
	if err != nil {
		return err
	}
	if !setting.MergePullRequest {
		return nil
	}

	client, err := s.SCM.Client(repo.SCM)
	if err != nil {
		return err
	}
	user, err := s.RepoStore.Creator(repo)
	if err != nil {
		return err
	}
	changes, err := client.PullRequests().ListChanges(ctx, user, repo.FullName(), hook.Number)
	if err != nil {
		return err
	}

	reports, err := s.ReportStore.Finds(&core.Report{
		ReportID: repo.ReportID,
		Commit:   hook.Commit,
	})
	if err != nil {
		return err
	}

	for _, report := range reports {
		previous, err := s.ReportStore.Find(&core.Report{
			ReportID: repo.ReportID,
			Branch:   hook.Target,
			Type:     report.Type,
		})
		if err != nil {
			report.Branch = hook.Target
		} else {
			report, err = s.ReportService.MergeReport(previous, report, changes)
			if err != nil {
				return err
			}
		}
		// TODO: need to use transation to prevent error occur in the middle
		if err := s.ReportStore.Upload(report); err != nil {
			return err
		}
	}
	return nil
}
