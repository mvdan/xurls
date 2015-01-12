/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package xurls

import "regexp"

var Regexp = regexp.MustCompile(
	`(([^\s,'"<>\(\)]+:(//)?|(http|ftp|www)[^.,;:'"<>\(\)]*\.)[^\s,'"<>\(\)]*[^.,;:\s.,;:'"<>\(\)]|[^\s,'"<>\(\)]+\.(com|org|net|edu|info)(/([^\s'"<>\(\)]*[^.,;:\s'"<>\(\)])?)?)`)
