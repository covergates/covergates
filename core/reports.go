package core

type StatementHit struct {
	LineNumber int64
	Hits       int64
}

type CoverageReport interface {
	StatementCov() float32
	Files() []CoverageFile
}

type CoverageFile interface {
	StatementCov() float32
	StatementHits() []StatementHit
}
