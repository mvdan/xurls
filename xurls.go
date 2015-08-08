// Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package xurls

import "regexp"

//go:generate go run generate/tldsgen/main.go
//go:generate go run generate/regexgen/main.go

const (
	letter    = `\p{L}`
	number    = `\p{N}`
	iriChar   = letter + number
	pathChar  = iriChar + `/\-+_@|&=#~*%.,:;'?!` + `\p{Sc}\p{Sk}\p{So}`
	endChar   = iriChar + `/\-+_&~*%` + `\p{Sc}`
	octet     = `(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])`
	ipv4Addr  = `\b` + octet + `\.` + octet + `\.` + octet + `\.` + octet + `\b`
	ipv6Addr  = `([0-9a-fA-F]{1,4}:([0-9a-fA-F]{1,4}:([0-9a-fA-F]{1,4}:([0-9a-fA-F]{1,4}:([0-9a-fA-F]{1,4}:[0-9a-fA-F]{0,4}|:[0-9a-fA-F]{1,4})?|(:[0-9a-fA-F]{1,4}){0,2})|(:[0-9a-fA-F]{1,4}){0,3})|(:[0-9a-fA-F]{1,4}){0,4})|:(:[0-9a-fA-F]{1,4}){0,5})((:[0-9a-fA-F]{1,4}){2}|:(25[0-5]|(2[0-4]|1[0-9]|[1-9])?[0-9])(\.(25[0-5]|(2[0-4]|1[0-9]|[1-9])?[0-9])){3})|(([0-9a-fA-F]{1,4}:){1,6}|:):[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){7}:`
	ipAddr    = `(` + ipv4Addr + `|` + ipv6Addr + `)`
	iri       = `[` + iriChar + `]([` + iriChar + `\-]*[` + iriChar + `])?`
	domain    = `(` + iri + `\.)+`
	hostName  = `(` + domain + gtld + `|` + ipAddr + `)`
	wellParen = `\([` + pathChar + `]*(\([` + pathChar + `]*\)[` + pathChar + `]*)*\)`
	pathCont  = `([` + pathChar + `]*(` + wellParen + `|[` + endChar + `])+)+`
	path      = `(/|/` + pathCont + `?|\b|$)`
	port      = `(:[0-9]*)?`
	webURL    = hostName + port + path
	email     = `[a-zA-Z0-9._%\-+]+@` + hostName

	comScheme = `[a-zA-Z][a-zA-Z.\-+]*://`
	scheme    = `(` + comScheme + `|` + otherScheme + `)`
	strict    = `\b` + scheme + pathCont
	relaxed   = strict + `|` + webURL + `|` + email
)

var (
	// Relaxed matches all the urls it can find.
	Relaxed = regexp.MustCompile(relaxed)
	// Strict only matches urls with a scheme to avoid false positives.
	Strict = regexp.MustCompile(strict)
)

func init() {
	Relaxed.Longest()
	Strict.Longest()
}

// StrictMatching produces a regexp that matches urls like Strict but whose
// scheme matches the given regular expression.
func StrictMatching(schemeExp string) (*regexp.Regexp, error) {
	strictMatching := `\b(?i)(` + schemeExp + `)(?-i)` + pathCont
	re, err := regexp.Compile(strictMatching)
	if err != nil {
		return nil, err
	}
	re.Longest()
	return re, nil
}
