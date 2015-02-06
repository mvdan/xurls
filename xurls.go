/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package xurls

import (
	"regexp"
)

//go:generate go run tools/regexgen/main.go

var (
	WebUrl = regexp.MustCompile(webUrl)
	Email  = regexp.MustCompile(email)
	All    = regexp.MustCompile(all)
)
