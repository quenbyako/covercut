package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/quenbyako/covercut/cmd/sequence"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/cover"
)

type seqtype = sequence.Sequence[int]

func Action(c *cli.Context) error {
	if c.NArg() < 2 {
		panic("requires 2 arguments: path to package and cover file")
	}

	var profile []*cover.Profile
	var err error
	if path := os.Args[2]; path == "-" {
		profile, err = cover.ParseProfilesFromReader(os.Stdin)
	} else {
		profile, err = cover.ParseProfiles(path)
	}
	if err != nil {
		panic(err)
	}

	modpath, packagename := getRoot(os.Args[1])
	profile = sliceFilter(profile, func(p *cover.Profile) bool {
		if !strings.HasPrefix(p.FileName, packagename) {
			return false
		}

		p.FileName = filepath.Join(modpath, strings.TrimPrefix(p.FileName, packagename))

		return true
	})

	for i, p := range profile {
		profile[i] = filterFile(p)
	}

	if len(profile) == 0 {
		return cli.Exit("Cover profile doesn't contain any data", 0)
	}

	fmt.Printf("mode: %v\n", profile[0].Mode)
	for _, p := range profile {
		for _, b := range p.Blocks {
			file := packagename + strings.TrimPrefix(p.FileName, modpath)
			fmt.Printf("%v:%v.%v,%v.%v %v %v\n", file, b.StartLine, b.StartCol, b.EndLine, b.EndCol, b.NumStmt, b.Count)
		}
	}

	return nil
}

func filterFile(p *cover.Profile) *cover.Profile {
	ignored := getIgnoredBlocks(p.FileName)

	return &cover.Profile{
		FileName: p.FileName,
		Mode:     p.Mode,
		Blocks: sliceFilter(p.Blocks, func(b cover.ProfileBlock) bool {
			s := seqtype{Start: b.StartLine, Stop: b.EndLine}
			return len(s.CutMany(ignored...)) > 0
		}),
	}
}

func getIgnoredBlocks(path string) []seqtype {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	noCoverLines := []int{} // already sorted
	for _, cg := range f.Comments {
		for _, comment := range cg.List {
			if strings.HasPrefix(strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")), "cover:ignore") {
				noCoverLines = append(noCoverLines, fset.Position(comment.Slash).Line)
			}
		}
	}

	ignoreBlocks := []seqtype{}
	for _, decl := range f.Decls {
		ignoreBlocks = append(ignoreBlocks, getIgnoredDecl(fset, decl, noCoverLines)...)
	}

	return ignoreBlocks
}

func getIgnoredDecl(f *token.FileSet, d ast.Decl, noCoverLines []int) []seqtype {
	switch d := d.(type) {
	case *ast.FuncDecl:
		seq := seqtype{
			Start: f.Position(d.Body.Lbrace).Line,
			Stop:  f.Position(d.Body.Rbrace).Line,
		}

		if in(noCoverLines, seq.Start) {
			return []seqtype{seq}
		}

		return []seqtype{}

		//ignored := []seqtype{}
		//for _, stmt := range d.Body.List {
		//	ignored = append(ignored, getIgnoredStmt(f, stmt, noCoverLines)...)
		//}
		//return ignored

	//case *ast.GenDecl:
	//	if start := f.Position(d.Lparen).Line; in(noCoverLines, start) {
	//		end := f.Position(d.Rparen).Line
	//		return []seqtype{{start, end}}
	//	}
	//
	//	for _, spec := range d.Specs {
	//		_ = spec
	//	}

	default:
		return []seqtype{}
		panic(fmt.Sprintf("what??? %T", d))
	}
}

func getIgnoredSpec(f *token.FileSet, s ast.Spec, noCoverLines []int) []seqtype {
	if l := f.Position(s.Pos()).Line; in(noCoverLines, l) {
		return []seqtype{{l, f.Position(s.End()).Line}}
	}

	return []seqtype{}
}

func getRoot(fromPath string) (path, packageName string) {
	p, err := filepath.Abs(fromPath)
	if err != nil {
		panic(err)
	}

	const sep = string(filepath.Separator)
	items := strings.Split(strings.TrimPrefix(p, sep), sep)
	for len(items) > 0 {
		joined := sep + strings.Join(items, sep)

		gomodPath := joined + "/go.mod"
		if _, err := os.Stat(gomodPath); err != nil {
			items = items[:len(items)-1]
			continue
		}

		file, err := os.ReadFile(gomodPath)
		if err != nil {
			panic(err)
		}
		m, err := modfile.Parse(gomodPath, file, nil)
		if err != nil {
			panic(err)
		}
		return joined, m.Module.Mod.Path
	}

	return "", ""
}

func sliceFilter[S ~[]T, T any](s S, f func(T) bool) S {
	if len(s) == 0 {
		return s
	}

	n := 0
	for _, val := range s {
		if f(val) {
			s[n] = val
			n++
		}
	}

	return s[:n]
}

func in(items []int, p int) bool {
	_, ok := slices.BinarySearch(items, p)
	return ok
}
