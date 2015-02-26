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
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		matches := re.FindAllString(word, -1)
		for _, match := range matches {
			fmt.Println(match)
		}
	}
}
