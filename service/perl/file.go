package perl

import "github.com/code-devel-cover/CodeCover/core"

// FileCollection of Perl source codes
type FileCollection struct {
	collect map[string][]*core.File
}

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
	for i, c := range count.Statement {
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
		files = append(files, mergeFiles(collect))
	}
	return files
}
