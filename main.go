package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type LintError struct {
	Pos token.Position
	Err error
}

func (e *LintError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.Pos.Filename, e.Pos.Line, e.Err.Error())
}

func init() {
	log.SetFlags(0) // do not log time
	log.SetOutput(os.Stderr)
}

func main() {
	var options Options
	flagParser := flags.NewParser(&options, flags.Default)
	if _, err := flagParser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if options.Args.Root == "" {
		options.Args.Root = "."
	}
	if info, err := os.Stat(options.Args.Root); err != nil || !info.IsDir() {
		log.Fatalf("invalid directory path: %v", options.Args.Root)
	}

	config := options.GetConfig()
	linter := NewLinter(options.Args.Root, config)

	var lintErrors []*LintError

	err := filepath.Walk(options.Args.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("[WARN] %v", err)
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		visitor, err := NewVisitor(path, linter)
		if err != nil {
			return err
		}
		errors := visitor.Walk()
		lintErrors = append(lintErrors, errors...)
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
