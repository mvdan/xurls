// Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/mvdan/xurls"
)

var (
	matching = flag.String("m", "", "")
	relaxed  = flag.Bool("r", false, "")
)

func init() {
	flag.Usage = func() {
		p := func(format string, a ...interface{}) {
			fmt.Fprintf(os.Stderr, format, a...)
		}
		p("Usage: xurls [-h] [files]\n\n")
		p("If no files are given, it reads from standard input.\n\n")
		p("   -m <regexp>   only match urls whose scheme matches a regexp\n")
		p("                    example: 'https?://|mailto:'\n")
		p("   -r            also match urls without a scheme (relaxed)\n")
	}
}

func scan(re *regexp.Regexp, r io.Reader) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		for _, match := range re.FindAllString(word, -1) {
			fmt.Println(match)
		}
	}
}

func main() {
	flag.Parse()
	if *relaxed && *matching != "" {
		fmt.Fprintf(os.Stderr, "-r and -m at the same time don't make much sense.\n")
		os.Exit(1)
	}
	re := xurls.Strict
	if *relaxed {
		re = xurls.Relaxed
	} else if *matching != "" {
		var err error
		if re, err = xurls.StrictMatchingScheme(*matching); err != nil {
			fmt.Fprintf(os.Stderr, "invalid regular expression '%s': %v\n", *matching, err)
			os.Exit(1)
		}
	}
	args := flag.Args()
	if len(args) == 0 {
		scan(re, os.Stdin)
	}
	for _, path := range args {
		if path == "-" {
			scan(re, os.Stdin)
			continue
		}
		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not open file '%s': %v\n", path, err)
			os.Exit(1)
		}
		scan(re, file)
		file.Close()
	}
}
