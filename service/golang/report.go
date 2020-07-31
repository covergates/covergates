package golang

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/covergates/covergates/core"
	"github.com/covergates/covergates/modules/util"
	log "github.com/sirupsen/logrus"
)

const mode = "mode: "

var errCoverageFormat = errors.New("error coverage report format")

type record struct {
	name     string
	fromLine int
	toLine   int
	count    int
}

type recordSlice []*record
type fileMap map[string]*fileRecord
type fileRecord struct {
	name   string
	hitMap map[int]int
}

// CoverageService of go language
type CoverageService struct{}

// Report of golang
func (s *CoverageService) Report(ctx context.Context, data io.Reader) (*core.CoverageReport, error) {
	scanner := bufio.NewScanner(data)
	records := make(recordSlice, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		} else if len(line) >= len(mode) && line[0:len(mode)] == mode {
			continue
		}
		r, err := newRecord(line)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records.coverage(), nil
}

// Find the golang coverage report
func (s *CoverageService) Find(ctx context.Context, path string) (string, error) {
	if !util.IsDir(path) {
		return path, nil
	}
	report := ""
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Base(path) == "coverage.out" {
			report = path
			return io.EOF
		}
		return nil
	})
	if err != nil && err != io.EOF {
		return "", err
	}
	if report == "" {
		return report, fmt.Errorf("go coverage report not found")
	}
	return report, nil
}

// Open reader of report with given path
func (s *CoverageService) Open(ctx context.Context, path string) (io.Reader, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%s is not file", path)
	}
	file, err := os.Open(path)
	defer file.Close()
	buf := &bytes.Buffer{}
	io.Copy(buf, file)
	return buf, nil
}

func newRecord(line string) (*record, error) {
	tokens := strings.Split(line, ":")
	if len(tokens) < 2 {
		return nil, errCoverageFormat
	}
	file := tokens[0]
	tokens = strings.Split(tokens[1], " ")
	if len(tokens) < 3 {
		return nil, errCoverageFormat
	}
	hits, err := strconv.Atoi(tokens[len(tokens)-1])
	if err != nil {
		return nil, errCoverageFormat
	}
	tokens = strings.Split(tokens[0], ",")
	if len(tokens) < 2 {
		return nil, errCoverageFormat
	}
	from := splitLine(tokens[0])
	to := splitLine(tokens[1])

	return &record{
		name:     file,
		fromLine: from,
		toLine:   to,
		count:    hits,
	}, nil
}

func splitLine(line string) int {
	tokens := strings.Split(line, ".")
	num, err := strconv.Atoi(tokens[0])
	if err != nil {
		log.Error(err)
	}
	return num
}

func (records recordSlice) coverage() *core.CoverageReport {
	m := make(fileMap)
	for _, record := range records {
		m.addRecord(record)
	}
	return m.toCoverage()
}

func (m fileMap) addRecord(r *record) {
	file, ok := m[r.name]
	if !ok {
		file = &fileRecord{
			name:   r.name,
			hitMap: make(map[int]int),
		}
		m[r.name] = file
	}
	file.addRecord(r)
}

func (m fileMap) toCoverage() *core.CoverageReport {
	files := make([]*core.File, 0)
	for _, file := range m {
		files = append(files, file.toFile())
	}
	report := &core.CoverageReport{
		Files: files,
	}
	report.StatementCoverage = report.AvgStatementCoverage()
	return report
}

func (f *fileRecord) addRecord(r *record) {
	for i := r.fromLine; i <= r.toLine; i++ {
		f.hitMap[i] += r.count
	}
}

func (f *fileRecord) toFile() *core.File {
	lines := f.lines()
	hits := make([]*core.StatementHit, len(lines))
	sum := 0
	for i, line := range lines {
		hits[i] = &core.StatementHit{
			Hits:       f.hitMap[line],
			LineNumber: line,
		}
		if hits[i].Hits > 0 {
			sum++
		}
	}
	avg := 0.0
	if len(lines) > 0 {
		avg = float64(sum) / float64(len(lines))
	}
	return &core.File{
		Name:              f.name,
		StatementHits:     hits,
		StatementCoverage: avg,
	}
}

func (f *fileRecord) lines() []int {
	slice := make([]int, 0)
	for k := range f.hitMap {
		slice = append(slice, k)
	}
	sort.Ints(slice)
	return slice
}
