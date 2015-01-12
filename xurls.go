/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package xurls

import (
	"regexp"
)

var regexLink = regexp.MustCompile(
	`(([^\s'"<>\(\)]+:(//)?|(http|ftp|www)[^.]*\.)[^\s'"<>\(\)]*|[^\s'"<>\(\)]+\.(com|org|net|edu|info)(/[^\s'"<>\(\)]*)?)[^.,;:\s'"<>\(\)]`)

func FindString(s string) string {
	return regexLink.FindString(s)
}

func FindAllString(s string) []string {
	return regexLink.FindAllString(s, -1)
}
