package charts

import "github.com/covergates/covergates/core"

// ChartService provides charts
type ChartService struct{}

// CoverageDiffTreeMap of two coverage reports
func (service *ChartService) CoverageDiffTreeMap(old, new *core.CoverageReport) core.Chart {
	return NewCoverageDiffTreeMap(old, new)
}
