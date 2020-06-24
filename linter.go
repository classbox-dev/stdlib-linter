package main

import (
	"os"
	"path/filepath"
	"strings"
)

type sentinel struct{}

type Set map[string]sentinel

func (s Set) Contains(item string) bool {
	_, ok := s[item]
	return ok
}

type Linter struct {
	Root                   string
	allowedPackages        Set
	allowedPackagePrefixes []string
	bannedIds              map[string]Set
	bannedCalls            map[string]Set
	goroutinesBanned       bool
}

func NewLinter(root string, config *Config) *Linter {
	linter := new(Linter)
	linter.Root = root
	linter.goroutinesBanned = config.GoroutinesBanned

	linter.allowedPackages = Set{}
	for _, p := range config.Packages {
		linter.allowedPackages[p] = sentinel{}
	}
	linter.allowedPackagePrefixes = config.PackagePrefixes
	linter.bannedIds = map[string]Set{}
	for pkg, ids := range config.BannedIds {
		idsSet := Set{}
		for _, id := range ids {
			idsSet[id] = sentinel{}
		}
		linter.bannedIds[pkg] = idsSet
	}
	linter.bannedCalls = map[string]Set{}
	for pkg, calls := range config.BannedCalls {
		callSet := Set{}
		for _, call := range calls {
			callSet[call] = sentinel{}
		}
		linter.bannedCalls[pkg] = callSet
	}
	return linter
}

func (linter *Linter) IsValidPackage(importLiteral string) bool {
	pkg := strings.Replace(importLiteral, `"`, "", -1)
	for _, prefix := range linter.allowedPackagePrefixes {
		if strings.HasPrefix(pkg, prefix) {
			return true
		}
	}
	return linter.allowedPackages.Contains(pkg)
}

func (linter *Linter) IsValidId(path string, idLiteral string) bool {
	bannedIds, ok := linter.bannedIds["*"]
	if ok && bannedIds.Contains(idLiteral) {
		return false
	}
	for _, pkg := range Subpaths(path) {
		bannedIds, ok := linter.bannedIds[pkg]
		if ok && bannedIds.Contains(idLiteral) {
			return false
		}
	}
	return true
}

func (linter *Linter) IsValidCall(path string, callLiteral string) bool {
	bannedCalls, ok := linter.bannedCalls["*"]
	if ok && bannedCalls.Contains(callLiteral) {
		return false
	}
	for _, pkg := range Subpaths(path) {
		bannedCalls, ok := linter.bannedCalls[pkg]
		if ok && bannedCalls.Contains(callLiteral) {
			return false
		}
	}
	return true
}

func (linter *Linter) AreGoroutinesBanned() bool {
	return linter.goroutinesBanned
}

func Subpaths(path string) []string {
	parts := strings.Split(path, string(os.PathSeparator))
	sp := make([]string, 0)
	for i := 1; i <= len(parts); i++ {
		sp = append(sp, filepath.Join(parts[:i]...))
	}
	return sp
}
