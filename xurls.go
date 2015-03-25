/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package xurls

import (
	"regexp"
)

//go:generate go run tools/tldsgen/main.go
//go:generate go run tools/regexgen/main.go

const (
	letters    = "a-zA-Z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF"
	iriChar    = letters + `0-9`
	ipv4Addr   = `(25[0-5]|2[0-4][0-9]|[0-1][0-9]{2}|[1-9][0-9]|[1-9])\.(25[0-5]|2[0-4][0-9]|[0-1][0-9]{2}|[1-9][0-9]|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1][0-9]{2}|[1-9][0-9]|[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1][0-9]{2}|[1-9][0-9]|[0-9])`
	ipv6Addr   = `([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])`
	ipAddr     = `(` + ipv4Addr + `|` + ipv6Addr + `)`
	iri        = `[` + iriChar + `]([` + iriChar + `\-]{0,61}[` + iriChar + `]){0,1}`
	hostName   = `(` + iri + `\.)+` + gtld
	domainName = `(` + hostName + `|` + ipAddr + `|localhost)`
	webURL     = `((https?://([a-zA-Z0-9$-_.+!*'(),;?&=]{1,64}(:[a-zA-Z0-9$\-_.+!*'(),;?&=]{1,25})?\@)?)?(` + domainName + `)(:\d{1,5})?)(/([` + iriChar + `;/?:@&=#~\-.+!*'(),_])*)?`
	email      = `[a-zA-Z0-9._%\-+]{1,256}@` + domainName
	all        = `(` + webURL + `|` + email + `)`
)

// Regex expressions that match various kinds of urls and addresses
var (
	WebURL = regexp.MustCompile(webURL)
	Email  = regexp.MustCompile(email)
	All    = regexp.MustCompile(all)
)
