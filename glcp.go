package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

var (
	write = flag.Bool("w", false, "write result to (source) file instead of stdout")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "No package name")
		return
	}

	for _, pkgname := range args {
		files, err := getPkgFiles(pkgname)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		err = addCommentsToFiles(*write, files...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
}

func addCommentsToFiles(write bool, filenames ...string) error {
	fset := token.NewFileSet()
	files := make(map[string][]byte)
	for _, filename := range filenames {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		files[filename] = src
	}

	for filename, src := range files {
		f, err := addComments(fset, filename, src)
		if err != nil {
			return err
		}

		if write {
			newFile, err := os.OpenFile(filename, os.O_WRONLY, os.ModeAppend)
			if err != nil {
				return err
			}
			err = format.Node(newFile, fset, f)
			if err != nil {
				return err
			}

			newFile.Close()
		} else {
			err = format.Node(os.Stdout, fset, f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func addComments(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	astutil.Apply(f, astutil.ApplyFunc(func(c *astutil.Cursor) bool {
		switch val := c.Node().(type) {
		case *ast.FuncDecl:
			if val.Name.IsExported() {
				val.Doc = addDoc(f, val.Doc, val.Name.Name, val.Pos())
			}
		case *ast.GenDecl:
			var name string
			if val.Tok == token.CONST || val.Tok == token.VAR {
				name = findExportedValueSpecName(val.Specs)
			} else if val.Tok == token.TYPE {
				name = findExportedTypeSpecName(val.Specs)
			}
			if name != "" {
				val.Doc = addDoc(f, val.Doc, name, val.Pos())
			}
		}
		return true
	}), nil)

	sort.Sort(ByPos(f.Comments))
	return f, nil
}

func addDoc(file *ast.File, doc *ast.CommentGroup, name string, pos token.Pos) *ast.CommentGroup {
	if doc == nil || len(doc.List) == 0 {
		doc = &ast.CommentGroup{
			List: []*ast.Comment{
				&ast.Comment{
					Slash: pos - 1,
					Text:  fmt.Sprintf("// %s needs a comment (THIS IS A PLACEHOLDER)\n", name),
				},
			},
		}
		file.Comments = append(file.Comments, doc)
	} else if !strings.HasPrefix(doc.List[0].Text, "// "+name) {
		doc.List = append([]*ast.Comment{&ast.Comment{
			Slash: doc.List[0].Pos() - 1,
			Text:  fmt.Sprintf("// %s needs a comment (THIS IS A PLACEHOLDER)\n", name),
		}}, doc.List...)
	}

	return doc
}

// ByPos is slice of CommentGroup for sorting
type ByPos []*ast.CommentGroup

func (s ByPos) Len() int {
	return len(s)
}

func (s ByPos) Less(i, j int) bool {
	return s[i].List[0].Pos() < s[j].List[0].Pos()
}

func (s ByPos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Some Spec doesn't have comment on it
func head(file *ast.File, pos token.Pos) token.Pos {
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			if comment.End() == pos-1 {
				return head(file, comment.Pos())
			}
		}
	}
	return pos
}

func findExportedValueSpecName(specs []ast.Spec) string {
	for _, spec := range specs {
		s, ok := spec.(*ast.ValueSpec)
		if ok {
			for _, name := range s.Names {
				if name.IsExported() {
					if len(specs) == 1 {
						return name.Name
					}
					return "This block"
				}
			}
		}
	}
	return ""
}

func findExportedTypeSpecName(specs []ast.Spec) string {
	for _, spec := range specs {
		s, ok := spec.(*ast.TypeSpec)
		if ok && s.Name.IsExported() {
			return s.Name.Name
		}
	}
	return ""
}
