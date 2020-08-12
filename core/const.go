package core

// ReportType of the coverage. Normally using programing language as a type
type ReportType string

const (
	// ReportPerl language
	ReportPerl ReportType = "perl"
	// ReportGo language
	ReportGo ReportType = "go"
	// ReportPython language
	ReportPython ReportType = "python"
	// ReportRuby language
	ReportRuby ReportType = "ruby"
	// ReportLCOV of lcov report
	ReportLCOV ReportType = "lcov"
)

// SCMProvider of Git service
type SCMProvider string

const (
	// Gitea SCM
	Gitea SCMProvider = "gitea"
	// Github SCM
	Github SCMProvider = "github"
)

// ReportUpdateAction action type when new report update
type ReportUpdateAction string

const (
	// ActionAutoMerge with previous report
	ActionAutoMerge ReportUpdateAction = "merge"
	// ActionAppend report to a new record
	ActionAppend ReportUpdateAction = "append"
)
