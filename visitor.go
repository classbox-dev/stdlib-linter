package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

type Visitor struct {
	pkg    string
	linter *Linter
	file   *ast.File
	fset   *token.FileSet
	stack  []ast.Node
	errors []*LintError
}

func NewVisitor(path string, linter *Linter) (*Visitor, error) {
	fset := token.NewFileSet() // positions are relative to fset
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	pkg, err := filepath.Rel(linter.Root, filepath.Dir(path))
	if err != nil {
		return nil, err
	}
	visitor := &Visitor{
		pkg:    pkg,
		linter: linter,
		file:   file,
		fset:   fset,
	}
	return visitor, nil
}

func (v *Visitor) Walk() []*LintError {
	ast.Walk(v, v.file)
	return v.errors
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		v.stack = v.stack[0 : len(v.stack)-1]
		return v
	}

	pos := v.fset.Position(node.Pos())

	v.stack = append(v.stack, node)

	switch n := node.(type) {

	case *ast.ImportSpec:
		if !v.linter.IsValidPackage(n.Path.Value) {
			err := &LintError{pos, fmt.Errorf("%s package is banned", n.Path.Value)}
			v.errors = append(v.errors, err)
		}

	case *ast.Ident:
		// <- identifier
		if !v.linter.IsValidId(v.pkg, n.Name) {
			err := &LintError{pos, fmt.Errorf("`%s` is not allowed in %s", n.Name, v.pkg)}
			v.errors = append(v.errors, err)
		}

	case *ast.CallExpr:
		ie, ok := n.Fun.(*ast.Ident)
		if !ok {
			break
		}
		// <- function call
		if !v.linter.IsValidCall(v.pkg, ie.Name) {
			err := &LintError{pos, fmt.Errorf("`%s()` call is not allowed in %s", ie.Name, v.pkg)}
			v.errors = append(v.errors, err)
		}
	case *ast.GoStmt:
		if v.linter.AreGoroutinesBanned() {
			err := &LintError{pos, errors.New("goroutines are not allowed")}
			v.errors = append(v.errors, err)
		}

	}
	return v
}
