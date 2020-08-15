package charts

import "github.com/covergates/covergates/core"

// ChartService provides charts
type ChartService struct{}

// CoverageDiffTreeMap of two coverage reports
func (service *ChartService) CoverageDiffTreeMap(old, new *core.Report) core.Chart {
	return NewCoverageDiffTreeMap(old, new)
}

// RepoCard of repository status
func (service *ChartService) RepoCard(repo *core.Repo, report *core.Report) core.Chart {
	return NewRepoCard(repo, report)
}
