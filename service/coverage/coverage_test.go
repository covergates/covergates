package coverage

import (
	"context"
	"reflect"
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
)

func TestTrimReportFileName(t *testing.T) {

	fileNames := []string{
		"blib/main.pl", "/path/to/trim/dir/test.pl",
	}
	expectNames := []string{"main.pl", "dir/test.pl"}
	filters := []string{
		"^blib/",
		"^/path/to/trim/",
		"(",
	}
	files := make([]*core.File, len(fileNames))

	for i, name := range fileNames {
		files[i] = &core.File{
			Name: name,
		}
	}

	report := &core.CoverageReport{
		Files: files,
	}

	service := &Service{}
	err := service.TrimFileNames(context.Background(), report, filters)
	if err != nil {
		t.Error(err)
		return
	}

	trimFiles := make([]string, len(fileNames))
	for i, file := range report.Files {
		trimFiles[i] = file.Name
	}

	if !reflect.DeepEqual(trimFiles, expectNames) {
		t.Logf("\nexpect: %v\nget: %v\n", expectNames, trimFiles)
		t.Fail()
	}
}
