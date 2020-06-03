package core

type File struct {
	Name              string
	StatementCoverage float32
	StatementHits     []*StatementHit
}

type StatementHit struct {
	LineNumber int
	Hits       int
}

func (h *StatementHit) Copy() *StatementHit {
	return &StatementHit{
		LineNumber: h.LineNumber,
		Hits:       h.Hits,
	}
}
