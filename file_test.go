package main

import (
	"reflect"
	"testing"

	"github.com/motomux/pretty"
)

func TestGetPkgFiles(t *testing.T) {
	testcases := map[string]struct {
		in       string
		outFiles []string
		outErr   error
	}{
		"get_package": {
			in:       "./testdata",
			outFiles: []string{"testdata/1.go", "testdata/2.go"},
			outErr:   nil,
		},
		"get_package_with_dot": {
			in:       "./testdata/.",
			outFiles: []string{"testdata/1.go", "testdata/2.go"},
			outErr:   nil,
		},
	}

	for k, c := range testcases {
		t.Run(k, func(t *testing.T) {
			files, err := getPkgFiles(c.in)
			if !reflect.DeepEqual(files, c.outFiles) {
				t.Errorf("output files doesn't match: diff=%v", pretty.Diff(files, c.outFiles))
			}
			if err != c.outErr {
				t.Errorf("output err doesn't match: actual=%v expected=%v", err, c.outErr)
			}
		})
	}
}
