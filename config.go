package main

type Config struct {
	Packages        []string            `yaml:"packages"`
	PackagePrefixes []string            `yaml:"package_prefixes"`
	BannedIds       map[string][]string `yaml:"banned_ids"`
	BannedCalls     map[string][]string `yaml:"banned_calls"`
}

var defaultConfig = Config{
	Packages: []string{
		"fmt",
		"math",
		"errors",
		"strings",
		"strconv",
		"math/rand",
		"math/bits",
		"github.com/cheekybits/genny/generic",
	},
	PackagePrefixes: []string{"hsecode.com/stdlib"},
	BannedIds:       map[string][]string{},
	BannedCalls:     map[string][]string{},
}
