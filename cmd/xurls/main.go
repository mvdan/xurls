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

	"mvdan.cc/xurls/v2"
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

func scanPath(re *regexp.Regexp, path string) error {
	r := os.Stdin
	if path != "-" {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}
	bufr := bufio.NewReader(r)
	for {
		line, err := bufr.ReadBytes('\n')
		for _, match := range re.FindAll(line, -1) {
			fmt.Printf("%s\n", match)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if *relaxed && *matching != "" {
		errExit(fmt.Errorf("-r and -m at the same time don't make much sense"))
	}
	var re *regexp.Regexp
	if *relaxed {
		re = xurls.Relaxed()
	} else if *matching != "" {
		var err error
		if re, err = xurls.StrictMatchingScheme(*matching); err != nil {
			errExit(err)
		}
	} else {
		re = xurls.Strict()
	}
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}
	for _, path := range args {
		if err := scanPath(re, path); err != nil {
			errExit(err)
		}
	}
}

func errExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
