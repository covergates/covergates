package perl

import "fmt"

type coverDB struct {
	Runs map[string]*coverRun `json:"runs"`
}

type coverRun struct {
	Name    string                 `json:"name"`
	Counts  map[string]*coverCount `json:"count"`
	Dir     string                 `json:"dir"`
	Digests map[string]string      `json:"digests"`
}

type coverCount struct {
	Statement []int `json:"statement"`
}

type coverDigest struct {
	Statement []int  `json:"statement"`
	File      string `json:"file"`
}

func updateSummary(s map[string]*coverCount, r *coverRun) error {
	for name, count := range r.Counts {
		sc, ok := s[name]
		if !ok {
			s[name] = newCoverCount(len(count.Statement))
			sc = s[name]
		}
		if err := sc.updateStatement(count.Statement); err != nil {
			return err
		}
	}
	return nil
}

func (db *coverDB) CountSummarize() (map[string]*coverCount, error) {
	if len(db.Runs) <= 0 {
		return nil, fmt.Errorf("No runs")
	}
	summary := make(map[string]*coverCount)
	for _, run := range db.Runs {
		if err := updateSummary(summary, run); err != nil {
			return nil, err
		}
	}
	return summary, nil
}

func newCoverCount(lines int) *coverCount {
	return &coverCount{
		Statement: make([]int, lines),
	}
}

func (c *coverCount) updateStatement(statement []int) error {
	if len(statement) != len(c.Statement) {
		return fmt.Errorf("Statement count mismatch")
	}
	for i, n := range statement {
		c.Statement[i] += n
	}
	return nil
}
