package filter

import (
	"regexp"
	"strings"
)

type filterType int

const (
	contains filterType = iota
	equal
	suffix
	prefix
	regex
	notContains
	notEqual
	notSuffix
	notPrefix
	notRegex
	unsupported
)

func parse(f string) (filterType, string) {
	patternFilterTypeMap := map[string]filterType{
		`^contains\((.+)\)$`:  contains,
		`^equal\((.+)\)$`:     equal,
		`^suffix\((.+)\)$`:    suffix,
		`^prefix\((.+)\)$`:    prefix,
		`^regex\((.+)\)$`:     regex,
		`^!contains\((.+)\)$`: notContains,
		`^!equal\((.+)\)$`:    notEqual,
		`^!suffix\((.+)\)$`:   notSuffix,
		`^!prefix\((.+)\)$`:   notPrefix,
		`^!regex\((.+)\)$`:    notRegex,
	}
	for p, t := range patternFilterTypeMap {
		r := regexp.MustCompile(p)
		m := r.FindAllStringSubmatch(f, -1)
		if len(m) > 0 {
			return t, m[0][1]
		}
	}
	return unsupported, f
}

// F filter function signature
type F func(string) bool

// Filter generate and return a filter function
func Filter(f string) F {
	ft, pattern := parse(f)
	m := map[filterType]F{
		contains:    func(t string) bool { return strings.Contains(pattern, t) },
		notContains: func(t string) bool { return !strings.Contains(pattern, t) },
		equal:       func(t string) bool { return pattern == t },
		notEqual:    func(t string) bool { return pattern != t },
		suffix:      func(t string) bool { return strings.HasSuffix(t, pattern) },
		notSuffix:   func(t string) bool { return !strings.HasSuffix(t, pattern) },
		prefix:      func(t string) bool { return strings.HasPrefix(t, pattern) },
		notPrefix:   func(t string) bool { return !strings.HasPrefix(t, pattern) },
		regex: func(t string) bool {
			r, e := regexp.Compile(pattern)
			if e != nil {
				return false
			}
			return r.MatchString(t)
		},
		notRegex: func(t string) bool {
			r, e := regexp.Compile(pattern)
			if e != nil {
				return false
			}
			return !r.MatchString(t)
		},
		unsupported: func(string) bool { return true },
	}
	return m[ft]
}
