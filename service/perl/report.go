package perl

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

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
const digiestFolder = "structure"

type DigestsMap map[string]*coverDigest
type CoverageService struct{}

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

func findDigests(files []*zip.File) (DigestsMap, error) {
	digests := make(DigestsMap)
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

func report(db *coverDB, digests DigestsMap) (*core.CoverageReport, error) {
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
