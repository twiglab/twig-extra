package remux

import (
	"regexp"
	"strings"

	"github.com/twiglab/twig"
)

type UrlRegex struct {
	Regex *regexp.Regexp
}

func (u *UrlRegex) Match(url string) (twig.UrlParams, bool) {
	params := u.Regex.SubexpNames()
	matches := u.Regex.FindAllStringSubmatch(url, -1)

	if matches == nil {
		return nil, false
	}

	result := make(twig.UrlParams)
	for i, n := range matches[0] {
		if len(params[i]) > 0 {
			result[params[i]] = n
		}
	}

	return result, true
}

type UrlRegexFunc func(string) *UrlRegex

func Full(regex string) *UrlRegex {
	return &UrlRegex{
		Regex: regexp.MustCompile(regex),
	}
}

func Pattern(pattern string) *UrlRegex {
	parts := strings.Split(pattern, "/")
	regex := "^"

	for i, part := range parts {
		if len(part) > 0 {
			if i > 0 {
				regex += "\\/"
			}

			// do we have special character?
			switch part[0] {
			case ':':
				groupName := "(?P<" + part[1:] + ">"
				regex += groupName + ".[^\\/]*)"
			case '*':
				//support named wildcards
				if len(part) > 1 {
					groupName := "(?P<" + part[1:] + ">"
					regex += groupName + ".*)"
				} else {
					regex += ".*"
				}
				if part[len(part)-1] == '/' {
					regex += "\\/"
				}
			default:
				regex += regexp.QuoteMeta(part)
			}
		}
	}

	if pattern[len(pattern)-1] == '/' {
		regex += "\\/$"
	} else {
		regex += "$"
	}

	return &UrlRegex{
		Regex: regexp.MustCompile(regex),
	}
}
