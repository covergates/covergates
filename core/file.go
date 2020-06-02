package core

type File struct {
	Name              string
	Path              string
	StatementCoverage float32
	StatementHits     []*StatementHit
}

type StatementHit struct {
	LineNumber int64
	Text       string
	Hits       int64
}
