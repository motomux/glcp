package main

import (
	"bytes"
	"go/format"
	"go/token"
	"testing"
)

func TestAddComments(t *testing.T) {
	testcases := map[string]struct {
		inSrc      []byte
		inFilename string
		outErr     error
		postSrc    string
	}{
		"func": {
			inSrc: []byte(`
			package hello
			
			func HelloFunc() {
				fmt.Println("hello")
			}
			`),
			inFilename: "Hello.go",
			outErr:     nil,
			postSrc: `package hello

// HelloFunc needs a comment (THIS IS A PLACEHOLDER)
func HelloFunc() {
	fmt.Println("hello")
}
`,
		},

		"var": {
			inSrc: []byte(`
			package hello

			var HelloVar = "hello"
			`),
			inFilename: "Hello.go",
			outErr:     nil,
			postSrc: `package hello

// HelloVar needs a comment (THIS IS A PLACEHOLDER)
var HelloVar = "hello"
`,
		},

		"const": {
			inSrc: []byte(`
			package hello

			const HelloConst = "hello"
			`),
			inFilename: "Hello.go",
			outErr:     nil,
			postSrc: `package hello

// HelloConst needs a comment (THIS IS A PLACEHOLDER)
const HelloConst = "hello"
`,
		},

		"type": {
			inSrc: []byte(`
			package hello

			type HelloStruct struct{
				Str string
			}
			`),
			inFilename: "Hello.go",
			outErr:     nil,
			postSrc: `package hello

// HelloStruct needs a comment (THIS IS A PLACEHOLDER)
type HelloStruct struct {
	Str string
}
`,
		},

		"block": {
			inSrc: []byte(`
			package hello

			var (
				HelloVar = "hello"
				HelloVar2 = "hello2"
			)
			`),
			inFilename: "Hello.go",
			outErr:     nil,
			postSrc: `package hello

// This block needs a comment (THIS IS A PLACEHOLDER)
var (
	HelloVar  = "hello"
	HelloVar2 = "hello2"
)
`,
		},

		"unexported_block": {
			inSrc: []byte(`
			package hello

			var (
				helloVar = "hello"
				helloVar2 = "hello2"
			)
			`),
			inFilename: "Hello.go",
			outErr:     nil,
			postSrc: `package hello

var (
	helloVar  = "hello"
	helloVar2 = "hello2"
)
`,
		},

		"mixed_block": {
			inSrc: []byte(`
			package hello

			var (
				helloVar = "hello"
				HelloVar2 = "hello2"
			)
			`),
			inFilename: "Hello.go",
			outErr:     nil,
			postSrc: `package hello

// This block needs a comment (THIS IS A PLACEHOLDER)
var (
	helloVar  = "hello"
	HelloVar2 = "hello2"
)
`,
		},

		"single_block": {
			inSrc: []byte(`
			package hello

			var (
				HelloVar = "hello"
			)
			`),
			inFilename: "Hello.go",
			outErr:     nil,
			postSrc: `package hello

// HelloVar needs a comment (THIS IS A PLACEHOLDER)
var (
	HelloVar = "hello"
)
`,
		},
	}

	for k, c := range testcases {
		t.Run(k, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := addComments(fset, c.inFilename, c.inSrc)
			if err != c.outErr {
				t.Errorf("out error doesn't match: actual=%v expected=%v", err, c.outErr)
			}

			buf := bytes.NewBuffer(nil)
			err = format.Node(buf, fset, f)
			if err != nil {
				t.Errorf("Failed to format: err=%v", err)
			}
			if buf.String() != c.postSrc {
				t.Errorf("post src doesn't match: actual=%s expected=%s", buf.String(), c.postSrc)
			}
		})
	}
}
