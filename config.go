package main

type sentinel struct{}

var allowedPackages = map[string]sentinel{
	"errors":    {},
	"fmt":       {},
	"strings":   {},
	"strconv":   {},
	"math":      {},
	"math/rand": {},
	"math/bits": {},
	"github.com/cheekybits/genny/generic": {},
}

var allowedPackagePrefixes = []string{
	"hsecode.com/stdlib",
}

var bannedIDs = map[string][]string{}

var bannedCalls = map[string][]string{}
