package core

// File holds the coverage information of a single file
type File struct {
	Name              string
	StatementCoverage float64
	StatementHits     []*StatementHit
}

// FileDiff defines the coverage differences of files
type FileDiff struct {
	File                  *File
	StatementCoverageDiff float64
	Removed               bool
}

// FileChange defines file status
type FileChange struct {
	Path    string
	Added   bool
	Renamed bool
	Deleted bool
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
