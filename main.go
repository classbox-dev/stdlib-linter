package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func parseRoot() (string, error) {
	flag.Parse()
	switch {
	case flag.NArg() == 0:
		return ".", nil
	case flag.NArg() == 1:
		path := flag.Args()[0]
		if info, err := os.Stat(path); err != nil || !info.IsDir() {
			return "", fmt.Errorf("invalid directory path: %v", path)
		}
		return path, nil
	default:
		return "", errors.New("single positional argument is required")
	}
}

type LintError struct {
	Pos token.Position
	Err error
}

func (e *LintError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.Pos.Filename, e.Pos.Line, e.Err.Error())
}

func processFile(path string) ([]LintError, error) {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	pkg := filepath.Base(filepath.Dir(path))
	v := Visitor{fset: fset, pkg: pkg}
	ast.Walk(&v, f)
	return v.Errs, nil
}

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	flag.Usage = usage
}

func usage() {
	log.Println("usage: stdlib-linter [path]")
	flag.PrintDefaults()
}

func main() {
	root, err := parseRoot()
	if err != nil {
		log.Fatal(err)
	}

	var lintErrors []LintError

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("[WARN] %v", err)
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		es, err := processFile(path)
		if err != nil {
			return err
		}
		lintErrors = append(lintErrors, es...)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, err := range lintErrors {
		log.Printf("[ERR] %s", err.Error())
	}

	if len(lintErrors) > 0 {
		os.Exit(1)
	}
}
