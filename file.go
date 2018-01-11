package main

import (
	"go/build"
	"path/filepath"
)

func getPkgFiles(pkgname string) ([]string, error) {
	pkg, err := build.Import(pkgname, ".", 0)
	if err != nil {
		return nil, err
	}

	var files []string
	files = append(files, pkg.GoFiles...)
	if pkg.Dir != "." {
		for i, f := range files {
			files[i] = filepath.Join(pkg.Dir, f)
		}
	}
	return files, nil
}
