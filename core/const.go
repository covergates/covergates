package core

type ReportType string

const (
	ReportPerl ReportType = "perl"
)

type SCMProvider string

const (
	Gitea  SCMProvider = "gitea"
	Github SCMProvider = "github"
)
