package statix

import (
	"regexp"
	"strings"
)

// Pattern is a wrapper for regular expressions.
type Pattern struct {
	Regexp string
}

// NewPattern returns a Pattern based on a given `regexp`.
func NewPattern(regexp string) Pattern {
	return Pattern{
		Regexp: regexp,
	}
}

// NewExtensionPattern returns a Pattern with a regexp matching
// all the filenames with the extensions given in parameter.
func NewExtensionPattern(extensions ...string) Pattern {
	exts := []string{}
	for _, ext := range extensions {
		exts = append(exts, "(\\."+ext+"$)")
	}

	return Pattern{
		Regexp: strings.Join(exts, "|"),
	}
}

// Match tests if the `s` string matches the Pattern.
func (p Pattern) Match(s string) bool {
	res, err := regexp.MatchString(p.Regexp, s)
	return res && (err == nil)
}
