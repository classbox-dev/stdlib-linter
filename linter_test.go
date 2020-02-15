package main_test

import (
	linter "github.com/mkuznets/stdlib-linter"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSubpaths(t *testing.T) {
	path := filepath.Join("a", "b", "c")
	expected := []string{
		filepath.Join("a"),
		filepath.Join("a", "b"),
		filepath.Join("a", "b", "c"),
	}
	result := linter.Subpaths(path)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Subpaths(%v) = %v, expected %v", path, result, expected)
	}

	result = linter.Subpaths("a")
	if !reflect.DeepEqual(result, []string{"a"}) {
		t.Fatalf("Subpaths(`a`) = %v, expected [a]", result)
	}
}
