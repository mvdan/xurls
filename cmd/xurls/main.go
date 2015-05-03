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
	matching = flag.String("m", "", "")
	relaxed  = flag.Bool("r", false, "")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `Usage: xurls [-h]`)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, `   -m <regexp>   only match urls whose scheme matches a regexp`)
		fmt.Fprintln(os.Stderr, `                    example: "https?://|mailto:"`)
		fmt.Fprintln(os.Stderr, `   -r            also match urls without a scheme (relaxed)`)
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
		var err error
		if re, err = xurls.StrictMatching(*matching); err != nil {
			fmt.Fprintln(os.Stderr, "invalid -m regular expression:", *matching)
			os.Exit(2)
		}
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		for _, match := range re.FindAllString(word, -1) {
			fmt.Println(match)
		}
	}
}
