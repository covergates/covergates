package core

import "io"

//go:generate mockgen -package mock -destination ../mock/chart_mock.go . ChartService,Chart

// ChartService provides charts
type ChartService interface {
	CoverageDiffTreeMap(old, new *CoverageReport) Chart
}

// Chart renders image to writer
type Chart interface {
	Render(w io.Writer) error
}
