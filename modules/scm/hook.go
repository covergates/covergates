package scm

import (
	"errors"
	"net/http"

	"github.com/covergates/covergates/config"
	"github.com/covergates/covergates/core"
	"github.com/drone/go-scm/scm"
)

var errWebhookNotSuport = errors.New("webhook not support")

type webhookService struct {
	config *config.Config
	client *scm.Client
	scm    core.SCMProvider
}

func (service *webhookService) Parse(req *http.Request) (core.HookEvent, error) {
	config := service.config
	hook, err := service.client.Webhooks.Parse(req, func(webhook scm.Webhook) (string, error) {
		return config.Server.Secret, nil
	})
	if err != nil {
		return nil, err
	}

	switch event := hook.(type) {
	case *scm.PullRequestHook:
		if event.Action == scm.ActionMerge {
			return &core.PullRequestHook{
				Number: event.PullRequest.Number,
				Merged: true,
				Commit: event.PullRequest.Sha,
				Source: event.PullRequest.Source,
				Target: event.PullRequest.Target,
			}, nil
		} else if service.scm == core.Gitea && event.Action == scm.ActionClose {
			return &core.PullRequestHook{
				Number: event.PullRequest.Number,
				Merged: true,
				Commit: event.PullRequest.Sha,
				Source: event.PullRequest.Source,
				Target: event.PullRequest.Target,
			}, nil
		}
	}
	return nil, errWebhookNotSuport
}

func (service *webhookService) IsWebhookNotSupport(err error) bool {
	return err == errWebhookNotSuport
}
