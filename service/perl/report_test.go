package perl

import (
	"archive/zip"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type Node struct {
	name    string
	entries []*Node // nil if the entry is a file
	mark    int
}

func walkTree(n *Node, path string, f func(path string, n *Node)) {
	f(path, n)
	for _, e := range n.entries {
		walkTree(e, filepath.Join(path, e.name), f)
	}
}

func makeTree(t *testing.T, tree *Node) {
	walkTree(tree, tree.name, func(path string, n *Node) {
		if n.entries == nil {
			fd, err := os.Create(path)
			if err != nil {
				t.Errorf("makeTree: %v", err)
				return
			}
			fd.Close()
		} else {
			os.Mkdir(path, 0770)
		}
	})
}

func unzip(file, folder string) error {
	z, err := zip.OpenReader(file)
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)
	if err := os.Chdir(folder); err != nil {
		return err
	}
	for _, file := range z.File {
		if file.FileInfo().IsDir() {
			if err := os.Mkdir(file.Name, os.ModePerm); err != nil {
				return err
			}
		} else {
			fo, err := file.Open()
			if err != nil {
				return err
			}
			data, err := ioutil.ReadAll(fo)
			fo.Close()
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(file.Name, data, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}

func TestPerlCoverageService(t *testing.T) {
	file := filepath.Join("testdata", "cover_db.zip")
	f, err := os.Open(file)
	if err != nil {
		t.Error(err)
		return
	}
	s := &CoverageService{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	report, err := s.Report(ctx, f)
	if err != nil {
		t.Error(err)
		return
	}
	if len(report.Files) <= 0 {
		t.Fail()
		return
	}
	if report.StatementCoverage <= 0 {
		t.Fail()
		return
	}
}

func TestPerlCoverageWithStringStatementHits(t *testing.T) {
	file := filepath.Join("testdata", "JSON_db.zip")
	f, err := os.Open(file)
	if err != nil {
		t.Error(err)
		return
	}
	s := &CoverageService{}
	report, err := s.Report(context.Background(), f)
	if err != nil {
		t.Error(err)
		return
	}
	if len(report.Files) <= 0 {
		t.Fail()
		return
	}
	if report.StatementCoverage <= 0 {
		t.Fail()
		return
	}
}

func TestFindReport(t *testing.T) {
	tree := &Node{
		"repository",
		[]*Node{
			{"dir", []*Node{
				{"cover_db", nil, 0},
			}, 0},
			{"cover_db", []*Node{
				{"cover.14", nil, 0},
			}, 0},
		},
		0,
	}
	tmpDir, err := ioutil.TempDir("./", "TestWalk")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	originDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originDir)
	makeTree(t, tree)
	service := &CoverageService{}
	report, err := service.Find(context.Background(), tree.name)
	if err != nil {
		t.Error(err)
		return
	}
	if report != filepath.Join("repository", "cover_db") {
		t.Log(report)
		t.Fail()
	}
}

func TestOpenReport(t *testing.T) {
	file := filepath.Join("testdata", "cover_db.zip")
	tmpDir, err := ioutil.TempDir("./", "TestOpenReport")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := unzip(file, tmpDir); err != nil {
		t.Fatal(err)
	}

	service := &CoverageService{}
	r, err := service.Open(context.Background(), tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	report, err := service.Report(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	if report.StatementCoverage <= 0 {
		t.Fail()
	}
	if len(report.Files) <= 0 {
		t.Fail()
	}
}
