package report

import (
	"fmt"

	"github.com/covergates/covergates/core"
	"github.com/gin-gonic/gin"
	"github.com/narqo/go-badge"
)

// HandleGetBadge for the report id
// @Summary get badge for the report id
// @Tags Report
// @Param id path string true "report id"
// @Param latest query bool false "get latest report in main branch"
// @Success 200 {object} string "badge svg"
// @Router /reports/{id}/badge [get]
func HandleGetBadge(
	reportStore core.ReportStore,
	repoStore core.RepoStore,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		reportID := c.Param("id")
		report, err := getLatest(reportStore, repoStore, reportID)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		data, err := badge.RenderBytes(
			"Covergates",
			fmt.Sprintf("%d%%", int(report.Coverage.StatementCoverage*100)),
			"#00838F",
		)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.Header("Cache-Control", "max-age=600")
		c.Data(200, "image/svg+xml", data)
	}
}
