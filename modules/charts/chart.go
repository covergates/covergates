package charts

import "github.com/code-devel-cover/CodeCover/core"

type ChartService struct{}

func (service *ChartService) CoverageDiffTreeMap(old, new *core.CoverageReport) core.Chart {
	return NewCoverageDiffTreeMap(old, new)
}
