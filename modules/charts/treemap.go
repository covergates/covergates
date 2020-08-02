package charts

import (
	"io"
	"sort"

	svgchart "github.com/blueworrybear/svg-charts"
	"github.com/covergates/covergates/core"
)

type diffFiles struct {
	names []string
	files map[string]*core.File
}

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
	names := make([]string, 0)
	for name := range c.newFiles {
		names = append(names, name)
	}
	diffFiles := &diffFiles{
		names: names,
		files: c.newFiles,
	}
	sort.Sort(diffFiles)
	for _, name := range diffFiles.names {
		newFile := c.newFiles[name]
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
		} else if diff < 0 {
			color = colorDecrease
		}
		if diff != 0 {
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
				Width:  600,
				Height: 400,
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

func (f *diffFiles) Len() int { return len(f.names) }
func (f *diffFiles) Less(i, j int) bool {
	return f.files[f.names[i]].Name < f.files[f.names[j]].Name
}
func (f *diffFiles) Swap(i, j int) {
	f.names[i], f.names[j] = f.names[j], f.names[i]
}
