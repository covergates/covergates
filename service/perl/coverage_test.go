package perl

import "testing"

func TestCoverSummary(t *testing.T) {
	count1 := &coverCount{
		Statement: []int{1, 0, 0, 0},
	}
	count2 := &coverCount{
		Statement: []int{0, 1, 0, 0},
	}
	run1 := &coverRun{
		Name: "/test",
		Counts: map[string]*coverCount{
			"test.pl": count1,
		},
	}
	run2 := &coverRun{
		Name: "/test",
		Counts: map[string]*coverCount{
			"test.pl": count2,
		},
	}
	db := &coverDB{
		Runs: map[string]*coverRun{
			"1": run1,
			"2": run2,
		},
	}
	summary, err := db.CountSummarize()
	if err != nil {
		t.Error(err)
		return
	}
	count, ok := summary["test.pl"]
	if !ok {
		t.Log("Cannot found test.pl")
		t.Fail()
		return
	}
	expectStatement := []int{1, 1, 0, 0}
	if len(count.Statement) != len(expectStatement) {
		t.Fail()
		return
	}
	for i, n := range count.Statement {
		if n != expectStatement[i] {
			t.Fail()
			return
		}
	}

}
