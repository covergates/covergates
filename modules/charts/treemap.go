package charts

import (
	"io"

	svgchart "github.com/blueworrybear/svg-charts"
	"github.com/code-devel-cover/CodeCover/core"
)

const (
	colorIncrease = "#16C36D"
	colorDecrease = "#C3161B"
	colorNoChange = "#ADC7C5"
)

// CoverageDiffTreeMap for two coverage reports
type CoverageDiffTreeMap struct {
	oldFiles map[string]*core.File
	newFiles map[string]*core.File
}

// NewCoverageDiffTreeMap with two coverage reports
func NewCoverageDiffTreeMap(old, new *core.CoverageReport) *CoverageDiffTreeMap {
	oldFiles := make(map[string]*core.File)
	newFiles := make(map[string]*core.File)
	for _, file := range old.Files {
		oldFiles[file.Name] = file
	}
	for _, file := range new.Files {
		newFiles[file.Name] = file
	}
	return &CoverageDiffTreeMap{
		oldFiles: oldFiles,
		newFiles: newFiles,
	}
}

// Render chart to writer
func (c *CoverageDiffTreeMap) Render(w io.Writer) error {
	colors := make([]string, 0)
	labels := make([]string, 0)
	data := make([]interface{}, 0)
	for name, newFile := range c.newFiles {
		oldCover := 0.0
		oldFile, ok := c.oldFiles[name]
		if ok {
			oldCover = oldFile.StatementCoverage
		}
		diff := newFile.StatementCoverage - oldCover
		color := colorNoChange
		label := ""
		if diff > 0 {
			color = colorIncrease
			label = name
		} else if diff < 0 {
			color = colorDecrease
			label = name
		}
		colors = append(colors, color)
		labels = append(labels, label)
		data = append(data, len(newFile.StatementHits))
	}
	svg, err := svgchart.NewSVGChart(
		svgchart.Options{
			Series: []*svgchart.SeriesOption{
				{
					Data:   svgchart.SeriesData(data...),
					Colors: colors,
				},
			},
			Chart: &svgchart.ChartOptions{
				Type:   "treemap",
				Width:  800,
				Height: 600,
			},
			Labels: labels,
			LabelOption: &svgchart.LabelOptions{
				FontSize: 12,
				Color:    "#FFFFFF",
			},
		},
	)
	if err != nil {
		return err
	}
	return svg.Render(w)
}
