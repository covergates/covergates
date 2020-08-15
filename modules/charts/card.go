package charts

import (
	"fmt"
	"io"

	svg "github.com/ajstarks/svgo"
	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/modules/charts/icons"
	"github.com/dustin/go-humanize"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	fileIcon = "files.svg"
	listIcon = "list.svg"
	gitIcon  = "git.svg"
	logoIcon = "logo.svg"
	cssStyle = `
	<style>
	.score-circle-ring {
		stroke: #3f8167;
		stroke-width: 6;
		fill: none;
		opacity: 0.2;
	}
	.score-circle {
		stroke: #3f8167;
		stroke-width: 6;
		fill: none;
		stroke-linecap: round;
		stroke-dasharray: 283;
	}
	</style>
	`
)

// RepoCard render svg repo status card
type RepoCard struct {
	core.Chart
	repo   *core.Repo
	report *core.Report
	canvas *svg.SVG
}

// NewRepoCard render
func NewRepoCard(repo *core.Repo, report *core.Report) *RepoCard {
	return &RepoCard{
		repo:   repo,
		report: report,
	}
}

func (c *RepoCard) start(w io.Writer) {
	c.canvas = svg.New(w)
	c.canvas.Start(495, 195)
	c.canvas.Writer.Write([]byte(cssStyle))
	c.canvas.ClipPath(`id="a1"`)
	c.canvas.Rect(0, 0, 500, 170)
	c.canvas.ClipEnd()
	c.canvas.Roundrect(1, 1, 493, 193, 5, 5, `fill="#3b726b"`)
	c.canvas.Roundrect(1, 1, 493, 193, 5, 5, `fill="#333333"`, `clip-path="url(#a1)"`)
}

func (c *RepoCard) end() {
	c.canvas.End()
}

func (c *RepoCard) title() {
	c.canvas.Text(15, 32, c.repo.FullName(), `font: 600 24px 'Segoe UI', sans-serif`, `fill="#80b9a2"`)
}

func (c *RepoCard) files() {
	c.canvas.Translate(15, 75)
	c.canvas.Group(`transform="scale(0.7) translate(0, -22)"`)
	c.canvas.Writer.Write(icons.MustAsset(fileIcon))
	c.canvas.Gend()
	c.canvas.Text(32, 0, "Total Files", `font: 600 18px 'Segoe UI', sans-serif`, `fill="#e6e6e6"`)
	p := message.NewPrinter(language.English)
	c.canvas.Text(200, 0, p.Sprintf("%d", len(c.report.Files)), `font: 600 18px 'Segoe UI', sans-serif`, `fill="#e6e6e6"`)
	c.canvas.Gend()
}

func (c *RepoCard) hits() {
	c.canvas.Translate(15, 115)
	c.canvas.Translate(-1, -16)
	c.canvas.Gtransform("scale(0.18)")
	c.canvas.Writer.Write(icons.MustAsset(listIcon))
	c.canvas.Gend()
	c.canvas.Gend()
	c.canvas.Text(32, 0, "Hit Lines", `font: 600 18px 'Segoe UI', sans-serif`, `fill="#e6e6e6"`)
	p := message.NewPrinter(language.English)
	c.canvas.Text(200, 0, p.Sprintf("%d", c.hitLines()), `font: 600 18px 'Segoe UI', sans-serif`, `fill="#e6e6e6"`)
	c.canvas.Gend()
}

func (c *RepoCard) hitLines() int {
	hits := 0
	for _, coverage := range c.report.Coverages {
		for _, file := range coverage.Files {
			for _, line := range file.StatementHits {
				if line.Hits > 0 {
					hits++
				}
			}
		}
	}
	return hits
}

func (c *RepoCard) build() {
	c.canvas.Translate(15, 155)
	c.canvas.Gtransform("scale(0.11) translate(-140, -160)")
	c.canvas.Writer.Write(icons.MustAsset(gitIcon))
	c.canvas.Gend()
	c.canvas.Text(32, 0, "Recent Build", `font: 600 18px 'Segoe UI', sans-serif`, `fill="#e6e6e6"`)
	c.canvas.Text(200, 0, fmt.Sprintf("%s, %s", c.repo.Branch, humanize.Time(c.report.CreatedAt)), `font: 600 18px 'Segoe UI', sans-serif`, `fill="#e6e6e6"`)
	c.canvas.Gend()
}

func (c *RepoCard) brand() {
	c.canvas.Translate(180, 185)
	c.canvas.Gtransform("scale(0.08) translate(-5, -180)")
	c.canvas.Writer.Write(icons.MustAsset(logoIcon))
	c.canvas.Gend()
	c.canvas.Text(30, 0, "covergates©2020", "font-size:12px")
	c.canvas.Gend()
}

func (c *RepoCard) score() {
	c.canvas.Translate(420, 75)
	c.canvas.Circle(0, 0, 45, `class="score-circle-ring"`)
	c.canvas.Circle(
		0, 0, 45,
		`class="score-circle"`,
		fmt.Sprintf("stroke-dashoffset: %d", int(45*2*3.14*(1.0-c.report.StatementCoverage()))),
	)
	c.canvas.Translate(-35, 12)
	score := int(c.report.StatementCoverage() * 100)
	pad := 0
	if score < 100 {
		pad = 10
	}
	c.canvas.Text(pad, 0, fmt.Sprintf("%d", score), `font: 600 36px 'Segoe UI', sans-serif`, `fill="#3f8167"`)
	c.canvas.Text(55, 0, "%", `font: 800 14px 'Segoe UI', sans-serif`, `fill="#3f8167"`)
	c.canvas.Gend()
	c.canvas.Gend()
}

// Render card to writer
func (c *RepoCard) Render(w io.Writer) error {
	c.start(w)
	c.title()
	c.files()
	c.hits()
	c.build()
	c.brand()
	c.score()
	c.end()
	return nil
}
