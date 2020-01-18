package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type Visitor struct {
	pkg   string
	fset  *token.FileSet
	stack []ast.Node
	Errs  []LintError
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		v.stack = v.stack[0 : len(v.stack)-1]
		return v
	}

	pos := v.fset.Position(node.Pos())
	idBlackList, _ := bannedIDs[v.pkg]
	callBlackList, _ := bannedCalls[v.pkg]

	v.stack = append(v.stack, node)

	switch n := node.(type) {
	case *ast.BasicLit:
		if len(v.stack) < 2 {
			break
		}
		outer := v.stack[len(v.stack)-2]
		if _, ok := outer.(*ast.ImportSpec); !ok {
			break
		}
		// <- import literal
		if !isPackageAllowed(n.Value) {
			err := LintError{pos, fmt.Errorf("%s package is banned", n.Value)}
			v.Errs = append(v.Errs, err)
		}
	case *ast.Ident:
		// <- identifier
		for _, id := range idBlackList {
			if n.Name == id {
				err := LintError{pos, fmt.Errorf("`%s` is not allowed in %s", n.Name, v.pkg)}
				v.Errs = append(v.Errs, err)
			}
		}
	case *ast.CallExpr:
		ie, ok := n.Fun.(*ast.Ident)
		if !ok {
			break
		}
		// <- function call
		for _, call := range callBlackList {
			if ie.Name == call {
				err := LintError{pos, fmt.Errorf("`%s()` call is not allowed in %s", ie.Name, v.pkg)}
				v.Errs = append(v.Errs, err)
			}
		}
	}
	return v
}

func isPackageAllowed(importLit string) bool {
	pkg := strings.Replace(importLit, `"`, "", -1)
	for _, prefix := range allowedPackagePrefixes {
		if strings.HasPrefix(pkg, prefix) {
			return true
		}
	}
	_, ok := allowedPackages[pkg]
	return ok
}
