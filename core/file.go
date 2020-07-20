package core

// File holds the coverage information of a single file
type File struct {
	Name              string
	StatementCoverage float64
	StatementHits     []*StatementHit
}

// StatementHit records hit count for a single line
type StatementHit struct {
	LineNumber int
	Hits       int
}

// Copy to a new StatementHit object
func (h *StatementHit) Copy() *StatementHit {
	return &StatementHit{
		LineNumber: h.LineNumber,
		Hits:       h.Hits,
	}
}
