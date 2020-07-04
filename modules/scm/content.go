package scm

import (
	"context"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/drone/go-scm/scm"
)

type contentService struct {
	client *scm.Client
	scm    core.SCMProvider
}

func (service *contentService) ListAllFiles(
	ctx context.Context,
	user *core.User,
	repo, ref string,
) ([]string, error) {
	client := service.client
	ctx = withUser(ctx, service.scm, user)

	var path string
	paths := []string{""}
	files := make([]string, 0)
	seen := make(map[string]bool)
	for len(paths) > 0 {
		path, paths = paths[0], paths[1:len(paths)]
		contents, _, err := client.Contents.List(ctx, repo, path, ref, scm.ListOptions{})
		if err != nil {
			return []string{}, err
		}
		for _, content := range contents {
			if content.Kind == scm.ContentKindDirectory && !seen[content.Path] {
				seen[content.Path] = true
				paths = append(paths, content.Path)
			} else if content.Kind == scm.ContentKindFile {
				files = append(files, content.Path)
			}
		}
	}
	return files, nil
}
