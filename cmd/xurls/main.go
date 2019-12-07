// Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"mvdan.cc/xurls/v2"
)

var (
	matching = flag.String("m", "", "")
	relaxed  = flag.Bool("r", false, "")
	fix      = flag.Bool("fix", false, "")
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
		p("   -fix          overwrite urls that redirect\n")
	}
}

func scanPath(re *regexp.Regexp, path string) error {
	f := os.Stdin
	if path != "-" {
		var err error
		f, err = os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	bufr := bufio.NewReader(f)
	var fixedBuf bytes.Buffer
	anyFixed := false
	var broken []string
	for {
		line, err := bufr.ReadBytes('\n')
		offset := 0
		for _, pair := range re.FindAllIndex(line, -1) {
			// The indexes are based on the original line.
			pair[0] += offset
			pair[1] += offset
			match := line[pair[0]:pair[1]]
			if !*fix {
				fmt.Printf("%s\n", match)
				continue
			}
			u, err := url.Parse(string(match))
			if err != nil {
				continue
			}
			fixed := u.String()
			switch u.Scheme {
			case "http", "https":
				// See if the URL redirects somewhere.
				client := &http.Client{
					Timeout: 10 * time.Second,
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						if len(via) >= 10 {
							return errors.New("stopped after 10 redirects")
						}
						// Keep the fragment around.
						req.URL.Fragment = u.Fragment
						fixed = req.URL.String()
						return nil
					},
				}
				resp, err := client.Get(fixed)
				if err != nil {
					continue
				}
				if resp.StatusCode >= 400 {
					broken = append(broken, string(match))
				}
				resp.Body.Close()
			}
			if fixed != string(match) {
				// Replace the url, and update the offset.
				newLine := line[:pair[0]]
				newLine = append(newLine, fixed...)
				newLine = append(newLine, line[pair[1]:]...)
				offset += len(newLine) - len(line)
				line = newLine
				anyFixed = true
			}
		}
		if *fix {
			if path == "-" {
				os.Stdout.Write(line)
			} else {
				fixedBuf.Write(line)
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	if anyFixed && path != "-" {
		f.Close()
		// Overwrite the file, if we weren't reading stdin. Report its
		// path too.
		fmt.Println(path)
		if err := ioutil.WriteFile(path, fixedBuf.Bytes(), 0666); err != nil {
			return err
		}
	}
	if len(broken) > 0 {
		return fmt.Errorf("found %d broken urls in %q:\n%s", len(broken),
			path, strings.Join(broken, "\n"))
	}
	return nil
}

func main() { os.Exit(main1()) }

func main1() int {
	flag.Parse()
	if *relaxed && *matching != "" {
		fmt.Fprintln(os.Stderr, "-r and -m at the same time don't make much sense")
		return 1
	}
	var re *regexp.Regexp
	if *relaxed {
		re = xurls.Relaxed()
	} else if *matching != "" {
		var err error
		if re, err = xurls.StrictMatchingScheme(*matching); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
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
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	return 0
}
