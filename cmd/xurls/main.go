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
	relaxed = flag.Bool("r", false, "also match urls without scheme (relaxed)")
	matching = flag.String("m", "", "only match urls whose scheme matches a regexp (e.g. `https?://|mailto:`)")
)

func main() {
	flag.Parse()
	if *relaxed && *matching != "" {
		fmt.Fprintln(os.Stderr, "-r and -m at the same time don't make much sense.")
		os.Exit(1)
	}
	re := xurls.Strict
	if *relaxed {
		re = xurls.Relaxed
	} else if *matching != "" {
		re = xurls.StrictMatching(*matching)
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
