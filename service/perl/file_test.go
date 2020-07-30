package perl

import (
	"testing"

	"github.com/covergates/covergates/core"
)

func TestMergeFile(t *testing.T) {
	files := []*core.File{
		{
			Name: "test.pl",
			StatementHits: []*core.StatementHit{
				{
					LineNumber: 1,
					Hits:       1,
				},
				{
					LineNumber: 2,
					Hits:       0,
				},
			},
		},
		{
			Name: "test.pl",
			StatementHits: []*core.StatementHit{
				{
					LineNumber: 1,
					Hits:       0,
				},
				{
					LineNumber: 2,
					Hits:       1,
				},
			},
		},
	}

	file := mergeFiles(files)
	if file.Name != "test.pl" {
		t.Fail()
		return
	}
	if len(file.StatementHits) != 2 {
		t.Fail()
		return
	}
	expectHits := []int{1, 1}
	for i, hit := range file.StatementHits {
		if expectHits[i] != hit.Hits {
			t.Fail()
		}
	}
	if files[0].StatementHits[1].Hits != 0 {
		t.Log("file being changed")
		t.Fail()
	}

}
