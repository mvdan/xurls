/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/mvdan/xurls"
)

var (
	email = flag.Bool("e", false, "match e-mails instead of web urls")
)

func main() {
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		exp := xurls.WebUrl
		if *email {
			exp = xurls.Email
		}
		matches := exp.FindAllString(line, -1)
		if matches == nil {
			continue
		}
		for _, match := range matches {
			fmt.Println(match)
		}
	}
}
