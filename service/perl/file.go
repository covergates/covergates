package perl

import (
	"sort"

	"github.com/covergates/covergates/core"
	log "github.com/sirupsen/logrus"
)

// FileCollection of Perl source codes
type FileCollection struct {
	collect map[string][]*core.File
}

type statementSlice []*core.StatementHit

func (s statementSlice) Len() int           { return len(s) }
func (s statementSlice) Less(i, j int) bool { return s[i].LineNumber < s[j].LineNumber }
func (s statementSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func newStatementHits(digest *coverDigest) []*core.StatementHit {
	statements := make([]*core.StatementHit, len(digest.Statement))
	for i, line := range digest.Statement {
		statements[i] = &core.StatementHit{
			LineNumber: line,
		}
	}
	return statements
}

func newFile(name string, count *coverCount, digest *coverDigest) *core.File {
	statements := newStatementHits(digest)
	if len(count.Statement) > len(statements) {
		log.Warningf("%s statement count does match to digest, will ignore extra statements", name)
		log.Debug(count)
		log.Debug(digest)
	}
	for i, c := range count.Statement {
		if i >= len(statements) {
			continue
		}
		statements[i].Hits = c
	}
	return &core.File{
		Name:          name,
		StatementHits: statements,
	}
}

func newFileCollection() *FileCollection {
	return &FileCollection{
		collect: make(map[string][]*core.File),
	}
}

func (c *FileCollection) add(file *core.File) {
	files, ok := c.collect[file.Name]
	if !ok {
		files = make([]*core.File, 0)
	}
	c.collect[file.Name] = append(files, file)
}

func sumStatementCoverage(hits []*core.StatementHit) float64 {
	if len(hits) <= 0 {
		return 0.0
	}
	s := 0
	for _, hit := range hits {
		if hit.Hits > 0 {
			s++
		}
	}
	return float64(s) / float64(len(hits))
}

func mergeFiles(files []*core.File) *core.File {
	if len(files) <= 0 {
		return nil
	}
	seed := files[0]
	hits := make([]*core.StatementHit, len(seed.StatementHits))
	for i, hit := range seed.StatementHits {
		hits[i] = hit.Copy()
	}
	merged := &core.File{
		Name:          seed.Name,
		StatementHits: hits,
	}
	for _, file := range files[1:] {
		for i, hit := range file.StatementHits {
			merged.StatementHits[i].Hits += hit.Hits
		}
	}
	merged.StatementCoverage = sumStatementCoverage(merged.StatementHits)
	return merged
}

func (c *FileCollection) mergedFiles() []*core.File {
	files := make([]*core.File, 0)
	for _, collect := range c.collect {
		merged := mergeFiles(collect)
		mergeStatementHist(merged)
		files = append(files, merged)
	}
	return files
}

func mergeStatementHist(file *core.File) {
	hitMap := make(map[int]*core.StatementHit)
	for _, hit := range file.StatementHits {
		if h, ok := hitMap[hit.LineNumber]; ok {
			hit.Hits += h.Hits
		}
		hitMap[hit.LineNumber] = hit
	}
	statements := make(statementSlice, 0)
	for _, hit := range hitMap {
		statements = append(statements, hit)
	}
	sort.Sort(statements)
	file.StatementHits = statements
}
