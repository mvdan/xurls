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
	re := xurls.WebURL
	if *email {
		re = xurls.Email
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindAllString(line, -1)
		if matches == nil {
			continue
		}
		for _, match := range matches {
			fmt.Println(match)
		}
	}
}
