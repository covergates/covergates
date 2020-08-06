package models

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/covergates/covergates/core"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/gorm"
)

type MockCoverReport struct {
	Name string
}

type reportSlice []*core.Report

func (r reportSlice) Len() int { return len(r) }
func (r reportSlice) Less(i, j int) bool {
	a, b := r[i], r[j]
	k1, k2 := a.ReportID+a.Commit+string(a.Type), b.ReportID+b.Commit+string(b.Type)
	return k1 < k2
}
func (r reportSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func testExpectReports(t *testing.T, expect, results reportSlice) {
	now := time.Now()
	for _, report := range results {
		report.CreatedAt = now
		report.Coverage = nil
	}
	for _, report := range expect {
		report.CreatedAt = now
	}
	sort.Sort(results)
	sort.Sort(expect)
	if diff := cmp.Diff(expect, results); diff != "" {
		t.Fatal(diff)
	}
}

func TestReportStoreUpload(t *testing.T) {
	ctrl, service := getDatabaseService(t)
	defer ctrl.Finish()
	store := &ReportStore{
		DB: service,
	}
	m := &core.Report{
		ReportID: "1234",
		Commit:   "1234",
		Type:     core.ReportPerl,
	}
	if err := store.Upload(m); err != nil {
		t.Error(err)
		return
	}
	r := &core.CoverageReport{
		Files: []*core.File{
			{
				Name: "mock",
			},
		},
	}
	m.Coverage = r
	if err := store.Upload(m); err != nil {
		t.Error(err)
		return
	}
	session := store.DB.Session()
	first := &Report{}
	session = session.Where(&Report{ReportID: "1234", Commit: "1234"}).First(first)
	if err := session.Error; err != nil {
		t.Error(err)
		return
	}
	r = &core.CoverageReport{}
	if err := json.Unmarshal(first.Data, r); err != nil {
		t.Error(err)
		return
	}
	if len(r.Files) <= 0 || r.Files[0].Name != "mock" {
		t.Fail()
	}
	var count int
	session = session.Where(&Report{ReportID: "1234", Commit: "1234"}).Count(&count)
	if count != 1 {
		t.Fail()
	}
}
func TestReportUploadReference(t *testing.T) {
	const reportID = "TestReportUploadReference"
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := &ReportStore{DB: db}

	testReferenceCount := func(t *testing.T, db core.DatabaseService, ref *Reference, expect int) {
		var refs []*Reference
		if err := db.Session().Find(&refs, ref).Error; err != nil {
			t.Fatal(err)
		}
		if len(refs) != expect {
			t.Fatalf("new reference count %d not match %d", len(refs), expect)
		}
	}

	t.Run("should reuse reference when report update", func(t *testing.T) {
		report1 := &core.Report{
			Commit:    "abc",
			ReportID:  reportID,
			Type:      core.ReportPerl,
			Reference: "master",
		}
		if err := store.Upload(report1); err != nil {
			t.Fatal(err)
		}
		if err := store.Upload(report1); err != nil {
			t.Fatal(err)
		}
		testReferenceCount(t, db, &Reference{ReportID: reportID, Name: "master"}, 1)
	})

	t.Run("should reuse reference for new report", func(t *testing.T) {
		report2 := &core.Report{
			Commit:    "edf",
			ReportID:  reportID,
			Type:      core.ReportPerl,
			Reference: "master",
		}
		if err := store.Upload(report2); err != nil {
			t.Fatal(err)
		}
		testReferenceCount(t, db, &Reference{ReportID: reportID, Name: "master"}, 1)
		ref := &Reference{ReportID: reportID, Name: "master"}
		db.Session().Preload("Reports").First(ref, ref)
		if len(ref.Reports) != 2 {
			t.Fatal("should have 2 reports relate to master")
		}
	})
}

func TestReportFind(t *testing.T) {
	ctrl, service := getDatabaseService(t)
	defer ctrl.Finish()
	store := &ReportStore{DB: service}
	id := "TestReportFind"
	reports := []*core.Report{
		{
			ReportID:  id + "1",
			Reference: "master",
			Commit:    "commit1",
			Type:      core.ReportPerl,
		},
		{
			ReportID:  id + "1",
			Reference: "master",
			Commit:    "commit2",
			Type:      core.ReportPerl,
		},
		{
			ReportID:  id + "2",
			Reference: "branch1",
			Commit:    "commit1",
			Type:      core.ReportGo,
		},
		{
			ReportID:  id + "2",
			Reference: "branch2",
			Commit:    "commit2",
			Type:      core.ReportGo,
		},
		{
			ReportID:  id + "2",
			Reference: "branch2",
			Commit:    "commit3",
			Type:      core.ReportGo,
		},
		{
			ReportID:  id + "2",
			Reference: "master",
			Commit:    "commit4",
			Type:      core.ReportGo,
		},
	}

	for _, report := range reports {
		if err := store.Upload(report); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("should find latest created", func(t *testing.T) {
		rst, err := store.Find(&core.Report{
			ReportID: id + "1",
		})
		if err != nil {
			t.Fatal(err)
		}
		if rst.Commit != "commit2" {
			t.Fail()
		}
	})

	t.Run("should find with reference", func(t *testing.T) {
		queries := []*core.Report{
			{ReportID: id + "2", Reference: "branch1"},
			{ReportID: id + "2", Reference: "branch2"},
		}
		expects := []*core.Report{
			{Commit: "commit1", Reference: "branch1"},
			{Commit: "commit3", Reference: "branch2"},
		}

		for i, query := range queries {
			rst, err := store.Find(query)
			if err != nil {
				t.Fatal(err)
			}
			expect := expects[i]
			if rst.Commit != expect.Commit || rst.Reference != expect.Reference {
				t.Fail()
			}
		}
	})

	t.Run("should not found report with reference and empty report id", func(t *testing.T) {
		rst, err := store.Find(&core.Report{Reference: "master"})
		if err == nil || !gorm.IsRecordNotFoundError(err) {
			t.Log(rst)
			t.Fail()
		}
	})

	t.Run("should return error for non existing reference", func(t *testing.T) {
		rst, err := store.Find(&core.Report{ReportID: id + "2", Reference: "fake-branch"})
		if err == nil || !gorm.IsRecordNotFoundError(err) {
			t.Log(rst)
			t.Fail()
		}
	})
}

func TestReportFinds(t *testing.T) {
	queryString := func(query *core.Report) string {
		fields := make([]string, 0)
		if query.Commit != "" {
			fields = append(fields, "commit="+query.Commit)
		}
		if query.Reference != "" {
			fields = append(fields, "reference="+query.Reference)
		}
		if query.ReportID != "" {
			fields = append(fields, "report_id="+query.ReportID)
		}
		if string(query.Type) != "" {
			fields = append(fields, "type ="+string(query.Type))
		}
		return strings.Join(fields, ",")
	}

	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := &ReportStore{DB: db}
	id := "TestReportFinds"
	reports := reportSlice{
		{
			ReportID: id + "1",
			Commit:   "commit1",
			Type:     core.ReportPerl,
		},
		{
			ReportID: id + "2",
			Commit:   "commit1",
			Type:     core.ReportPerl,
		},
		{
			ReportID: id + "2",
			Commit:   "commit1",
			Type:     core.ReportGo,
		},
		{
			ReportID: id + "2",
			Commit:   "commit2",
			Type:     core.ReportGo,
		},
		{
			ReportID:  id + "3",
			Commit:    "commit1",
			Reference: "master",
			Type:      core.ReportGo,
		},
		{
			ReportID:  id + "3",
			Commit:    "commit2",
			Reference: "master",
			Type:      core.ReportGo,
		},
		{
			ReportID:  id + "3",
			Commit:    "commit2",
			Reference: "master",
			Type:      core.ReportPerl,
		},
	}

	queries := reportSlice{
		{ReportID: id + "1"},
		{ReportID: id + "2", Commit: "commit1"},
		{ReportID: id + "3", Reference: "master"},
	}

	expects := []reportSlice{
		{
			{
				ReportID: id + "1",
				Commit:   "commit1",
				Type:     core.ReportPerl,
			},
		},
		{
			{
				ReportID: id + "2",
				Commit:   "commit1",
				Type:     core.ReportPerl,
			},
			{
				ReportID: id + "2",
				Commit:   "commit1",
				Type:     core.ReportGo,
			},
		},
		{
			{
				ReportID:  id + "3",
				Commit:    "commit1",
				Reference: "master",
				Type:      core.ReportGo,
			},
			{
				ReportID:  id + "3",
				Commit:    "commit2",
				Reference: "master",
				Type:      core.ReportGo,
			},
			{
				ReportID:  id + "3",
				Commit:    "commit2",
				Reference: "master",
				Type:      core.ReportPerl,
			},
		},
	}

	for _, report := range reports {
		if err := store.Upload(report); err != nil {
			t.Fatal(err)
		}
	}

	for i, query := range queries {
		t.Run(queryString(query), func(t *testing.T) {
			expect := expects[i]
			results, err := store.Finds(query)
			if err != nil {
				t.Fatal(err)
			}
			testExpectReports(t, expect, reportSlice(results))
		})
	}
}

func TestReportList(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	store := &ReportStore{DB: db}
	id := "TestReportList"
	reports := reportSlice{
		{
			ReportID: id + "1",
			Commit:   "commit1",
			Type:     core.ReportPerl,
		},
		{
			ReportID: id + "1",
			Commit:   "commit1",
			Type:     core.ReportGo,
		},
		{
			ReportID:  id + "2",
			Commit:    "commit1",
			Type:      core.ReportGo,
			Reference: "master",
		},
		{
			ReportID:  id + "2",
			Commit:    "commit2",
			Type:      core.ReportGo,
			Reference: "master",
		},
	}

	queries := [][]string{
		{id + "1", "commit1"},
		{id + "2", "commit2"},
		{id + "2", "master"},
	}

	expectations := []reportSlice{
		{
			{
				ReportID: id + "1",
				Commit:   "commit1",
				Type:     core.ReportPerl,
			},
			{
				ReportID: id + "1",
				Commit:   "commit1",
				Type:     core.ReportGo,
			},
		},
		{
			{
				ReportID: id + "2",
				Commit:   "commit2",
				Type:     core.ReportGo,
			},
		},
		{
			{
				ReportID:  id + "2",
				Commit:    "commit1",
				Type:      core.ReportGo,
				Reference: "master",
			},
			{
				ReportID:  id + "2",
				Commit:    "commit2",
				Type:      core.ReportGo,
				Reference: "master",
			},
		},
	}

	for _, report := range reports {
		if err := store.Upload(report); err != nil {
			t.Fatal(err)
		}
	}
	base := 0
	for i, query := range queries[base:] {
		result, err := store.List(query[0], query[1])
		if err != nil {
			t.Fatal(err)
		}
		testExpectReports(t, expectations[i+base], reportSlice(result))
	}
}

func TestReportUploadFiles(t *testing.T) {
	ctrl, service := getDatabaseService(t)
	defer ctrl.Finish()
	store := &ReportStore{
		DB: service,
	}
	m := &core.Report{
		ReportID: "test_upload_files",
		Commit:   "test_upload_files",
		Files:    []string{"a", "b", "c"},
		Type:     core.ReportPerl,
	}
	if err := store.Upload(m); err != nil {
		t.Error(err)
		return
	}
	report, err := store.Find(&core.Report{
		ReportID: m.ReportID,
		Commit:   m.Commit,
	})
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(report.Files, m.Files) {
		t.Fail()
	}
}

func TestReportComment(t *testing.T) {
	ctrl, service := getDatabaseService(t)
	defer ctrl.Finish()
	store := &ReportStore{
		DB: service,
	}

	report := &core.Report{
		ReportID: "ABCD",
	}

	if err := store.CreateComment(report, &core.ReportComment{}); err == nil {
		t.Fail()
	}

	if err := store.CreateComment(report, &core.ReportComment{Comment: 1, Number: 1}); err != nil {
		t.Fatal(err)
	}
	comment, err := store.FindComment(report, 1)
	if err != nil {
		t.Fatal(err)
	}
	if comment.Comment != 1 {
		t.Fail()
	}
	if err := store.CreateComment(report, &core.ReportComment{Comment: 2, Number: 1}); err != nil {
		t.Fatal(err)
	}
	comment, err = store.FindComment(report, 1)
	if err != nil {
		t.Fatal(err)
	}
	if comment.Comment != 2 {
		t.Fail()
	}
	if _, err := store.FindComment(report, 123); err == nil {
		t.Fail()
	}
}
