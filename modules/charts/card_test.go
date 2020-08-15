package charts

import (
	"os"
	"testing"
	"time"

	"github.com/covergates/covergates/core"
)

func TestCard(t *testing.T) {
	card := NewRepoCard(
		&core.Repo{
			NameSpace: "covergates",
			Name:      "covergates",
			Branch:    "master",
		}, &core.Report{
			Files:     []string{"main.go", "core.go"},
			CreatedAt: time.Date(2020, time.July, 20, 12, 12, 40, 0, time.Local),
			Coverages: []*core.CoverageReport{
				{
					Files: []*core.File{
						{
							StatementCoverage: 0.8,
						},
					},
				},
			},
		},
	)
	file, err := os.Create("card.svg")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()
	card.Render(file)
}
