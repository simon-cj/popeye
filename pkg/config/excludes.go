package config

import (
	"regexp"
	"strings"

	"github.com/rs/zerolog"
)

// RxMarker indicate exclude flag is a regular expression.
const rxMarker = "rx:"

// RegExp defined regex to check if exclude is a regex or plain string.
var regExp = regexp.MustCompile(`\A` + rxMarker)

type (
	Exclusion struct {
		Name  string
		Codes []ID
	}

	// Exclude represents a collection of excludes items.
	// This can be a straight string match of regex using an rx: prefix.
	Exclusions []Exclusion

	// Excludes represents a set of resources that should be excluded
	// from the sanitizer.
	Excludes map[string]Exclusions
)

func init() {
	zerolog.SetGlobalLevel(zerolog.FatalLevel)
}

func newExcludes() Excludes {
	return Excludes{}
}

// ExcludeFQN checks if a given named resource should be excluded.
func (e Excludes) ExcludeFQN(section, fqn string) bool {
	excludes, ok := e[section]
	if !ok {
		return false
	}

	for _, exclude := range excludes {
		if exclude.Match(fqn) {
			return true
		}
	}

	return false
}

// ShouldExclude checks if a given named resource should be excluded.
func (e Excludes) ShouldExclude(section, fqn string, code ID) bool {
	// Not mentioned in config. Allow all
	excludes, ok := e[section]
	if !ok {
		return false
	}

	return excludes.Match(fqn, code)
}

// ShouldExclude checks if a given named should be excluded.
func (e Exclusions) Match(resource string, code ID) bool {
	for _, exclude := range e {
		if exclude.Match(resource) && hasCode(exclude.Codes, code) {
			return true
		}
	}

	return false
}

func (e Exclusion) Match(fqn string) bool {
	if !isRegex(e.Name) {
		return fqn == e.Name
	}

	return rxMatch(e.Name, fqn)
}

// Helpers...

func rxMatch(exp, name string) bool {
	rx := regexp.MustCompile(strings.Replace(exp, rxMarker, "", 1))
	b := rx.MatchString(name)
	return b
}

func isRegex(f string) bool {
	return regExp.MatchString(f)
}

func hasCode(codes []ID, code ID) bool {
	if len(codes) == 0 {
		return true
	}

	for _, c := range codes {
		if c == code {
			return true
		}
	}
	return false
}
