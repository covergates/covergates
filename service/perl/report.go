package perl

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/modules"
)

var errCoverDatabaseNotFound = errors.New("Coverage database not found")
var errDigestNoFound = errors.New("Digest not found")

const coverDBName = "cover.14"
const digiestFolder = "structure"

type PerlCoverageReportService struct {
	db      *coverDB
	files   map[string]*zip.File
	digests map[string]*coverDigest
}

func findDigests(files []*zip.File) (map[string]*coverDigest, error) {
	digests := make(map[string]*coverDigest)
	for _, file := range files {
		if file.FileInfo().IsDir() {
			continue
		}
		if filepath.Dir(file.Name) == digiestFolder {
			name := filepath.Base(file.Name)
			if filepath.Ext(name) == ".lock" {
				continue
			}
			digest := &coverDigest{}
			r, err := file.Open()
			defer r.Close()
			data, err := ioutil.ReadAll(r)
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(data, digest); err != nil {
				return nil, err
			}
			digests[name] = digest
		}
	}
	return digests, nil
}

func NewPerlCoverageReportService(data io.Reader) (*PerlCoverageReportService, error) {
	z, err := modules.NewZipReader(data)
	if err != nil {
		return nil, err
	}
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
	digests, err := findDigests(z.File)
	if err != nil {
		return nil, err
	}
	s := &PerlCoverageReportService{
		db:      db,
		digests: digests,
	}
	return s, nil
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

func (r *PerlCoverageReportService) Report(ctx context.Context) (*core.CoverageReport, error) {
	fileCollection := newFileCollection()
	for _, run := range r.db.Runs {
		for name, count := range run.Counts {
			key, ok := run.Digests[name]
			if !ok {
				return nil, errDigestNoFound
			}
			digest, ok := r.digests[key]
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
