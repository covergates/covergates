package core

// ReportType of the coverage. Normally using programing language as a type
type ReportType string

const (
	// ReportPerl language
	ReportPerl ReportType = "perl"
)

// SCMProvider of Git service
type SCMProvider string

const (
	// Gitea SCM
	Gitea SCMProvider = "gitea"
	// Github SCM
	Github SCMProvider = "github"
)
