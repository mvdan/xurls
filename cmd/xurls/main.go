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
	matching = flag.String("m", "", "only match urls whose scheme matches a regexp (e.g. `https?://|mailto:`)")
	relaxed  = flag.Bool("r", false, "also match urls without scheme (relaxed)")
)

func init() {
	flag.Usage = func() {
		p := func(args ...interface{}) {
			fmt.Fprintln(os.Stderr, args...)
		}
		p("Usage: xurls [-h]")
		p()
		p("   -m <regexp>  only match urls whose scheme matches a regexp")
		p("                   example: \"https?://|mailto:\"")
		p("   -r           also match urls without a scheme (relaxed)")
	}
}

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
