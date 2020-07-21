package perl

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/modules/archive"
	"github.com/code-devel-cover/CodeCover/modules/util"
)

var errCoverDatabaseNotFound = errors.New("Coverage database not found")
var errDigestNoFound = errors.New("Digest not found")

type errDigestFormat struct {
	msg string
}

func (err *errDigestFormat) Error() string {
	return fmt.Sprintf("digest format error: %s", err.msg)
}

const coverDBName = "cover.14"
const coverFolderName = "cover_db"
const digiestFolder = "structure"

type digestsMap map[string]*coverDigest

// CoverageService for Perl
type CoverageService struct{}

// Report coverage for Perl
func (r *CoverageService) Report(
	ctx context.Context,
	data io.Reader,
) (*core.CoverageReport, error) {

	z, err := archive.NewZipReader(data)
	if err != nil {
		return nil, err
	}
	db, err := findCoverDB(z)
	if err != nil {
		return nil, err
	}
	digests, err := findDigests(z.File)
	if err != nil {
		return nil, err
	}
	return report(db, digests)
}

// Find coverage report of Perl from given path
func (r *CoverageService) Find(ctx context.Context, path string) (string, error) {
	if !util.IsDir(path) {
		return "", fmt.Errorf("not found")
	}
	if filepath.Base(path) == coverFolderName {
		return path, nil
	}
	report := ""
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && filepath.Base(path) == coverFolderName {
			report = path
			return io.EOF
		}
		return nil
	})
	if err != nil && err != io.EOF {
		return "", err
	}
	return report, nil
}

// Open coverage report
func (r *CoverageService) Open(ctx context.Context, path string) (io.Reader, error) {
	if !util.IsDir(path) {
		return nil, fmt.Errorf("%s is not a folder", path)
	}
	root := strings.TrimRight(path, "/")
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	defer w.Close()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if len(path) <= len(root) {
			return nil
		}
		name := path[len(root)+1:]
		if info.IsDir() {
			if _, err := w.Create(name + "/"); err != nil {
				return err
			}
		} else {
			f, err := w.Create(name)
			if err != nil {
				return err
			}
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			if _, err := f.Write(data); err != nil {
				return err
			}
		}
		return nil
	})
	return buf, err
}

func findDigests(files []*zip.File) (digestsMap, error) {
	digests := make(digestsMap)
	for _, file := range files {
		if file.FileInfo().IsDir() {
			continue
		}
		if filepath.Dir(file.Name) == digiestFolder {
			name := filepath.Base(file.Name)
			if filepath.Ext(name) == ".lock" {
				continue
			}
			r, err := file.Open()
			defer r.Close()
			digest, err := unmarshalCoverDigest(r)
			if err != nil {
				return nil, err
			}
			digests[name] = digest
		}
	}
	return digests, nil
}

func unmarshalCoverDigest(r io.Reader) (*coverDigest, error) {
	var m map[string]interface{}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	fileData, ok := m["file"]
	if !ok {
		return nil, &errDigestFormat{msg: "file not found"}
	}
	file, ok := fileData.(string)
	if !ok {
		return nil, &errDigestFormat{msg: "file is not string"}
	}

	statementData, ok := m["statement"]
	if !ok {
		return nil, &errDigestFormat{msg: "statement not found"}
	}
	statementSlice, ok := statementData.([]interface{})
	if !ok {
		return nil, &errDigestFormat{msg: "statement is not array"}
	}

	statements, err := util.ToIntSlice(statementSlice)
	if err != nil {
		return nil, err
	}
	return &coverDigest{
		File:      file,
		Statement: statements,
	}, nil
}

func findCoverDB(z *zip.Reader) (*coverDB, error) {
	files := make(map[string]*zip.File)
	for _, file := range z.File {
		files[file.Name] = file
	}
	cover, ok := files[coverDBName]
	if !ok {
		return nil, errCoverDatabaseNotFound
	}
	content, err := cover.Open()
	defer content.Close()
	if err != nil {
		return nil, err
	}
	d, err := ioutil.ReadAll(content)
	if err != nil {
		return nil, err
	}
	db := &coverDB{}
	if err := json.Unmarshal(d, db); err != nil {
		return nil, err
	}
	return db, nil
}

func avgStatementCoverage(files []*core.File) float64 {
	if len(files) <= 0 {
		return 0.0
	}
	s := float64(0)
	for _, file := range files {
		s += file.StatementCoverage
	}
	return s / float64(len(files))
}

func report(db *coverDB, digests digestsMap) (*core.CoverageReport, error) {
	fileCollection := newFileCollection()
	for _, run := range db.Runs {
		for name, count := range run.Counts {
			key, ok := run.Digests[name]
			if !ok {
				return nil, errDigestNoFound
			}
			digest, ok := digests[key]
			if !ok {
				return nil, errDigestNoFound
			}
			fileCollection.add(newFile(name, count, digest))
		}
	}
	files := fileCollection.mergedFiles()
	report := &core.CoverageReport{
		StatementCoverage: avgStatementCoverage(files),
		Files:             files,
	}
	return report, nil
}
