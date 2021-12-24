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
	"sync/atomic"
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
	in := os.Stdin
	out := io.Writer(os.Stdout)
	var outBuf *bytes.Buffer
	if path != "-" {
		var err error
		in, err = os.Open(path)
		if err != nil {
			return err
		}
		if *fix {
			outBuf = new(bytes.Buffer)
			out = outBuf
		}
		defer in.Close()
	}

	// A maximum of 32 parallel requests.
	maxWeight := int64(32)
	seq := newSequencer(maxWeight, out, os.Stderr)

	scanner := bufio.NewScanner(in)

	// Doesn't need to be part of reporterState as order doesn't matter.
	var atomicFixedCount uint32

	for scanner.Scan() {
		line := scanner.Text() + "\n"
		matches := re.FindAllStringIndex(line, -1)
		if !*fix {
			for _, pair := range matches {
				match := line[pair[0]:pair[1]]
				fmt.Printf("%s\n", match)
			}
			continue
		}
		weight := int64(len(matches))
		if weight > maxWeight {
			weight = maxWeight
		}
		seq.Add(weight, func(r *reporter) error {
			offsetWithinLine := 0
			for _, pair := range re.FindAllStringIndex(line, -1) {
				// The indexes are based on the original line.
				pair[0] += offsetWithinLine
				pair[1] += offsetWithinLine
				match := line[pair[0]:pair[1]]
				origURL, err := url.Parse(match)
				if err != nil {
					continue
				}
				fixed := origURL.String()
				switch origURL.Scheme {
				case "http", "https":
					// See if the URL redirects somewhere.
					// Only apply a fix if the redirect chain is permanent.
					allPermanent := true
					client := &http.Client{
						Timeout: 10 * time.Second,
						CheckRedirect: func(req *http.Request, via []*http.Request) error {
							if len(via) >= 10 {
								return errors.New("stopped after 10 redirects")
							}
							switch req.Response.StatusCode {
							case http.StatusMovedPermanently, http.StatusPermanentRedirect:
							default:
								allPermanent = false
							}
							if allPermanent {
								// Inherit the fragment if empty.
								if req.URL.Fragment == "" {
									req.URL.Fragment = origURL.Fragment
								}
								fixed = req.URL.String()
							}
							return nil
						},
					}
					resp, err := client.Get(fixed)
					if err != nil {
						continue
					}
					if resp.StatusCode >= 400 {
						r.appendBroken(match)
					}
					resp.Body.Close()
				}
				if fixed != match {
					// Replace the url, and update offsetWithinLine.
					newLine := line[:pair[0]] + fixed + line[pair[1]:]
					offsetWithinLine += len(newLine) - len(line)
					line = newLine
					atomic.AddUint32(&atomicFixedCount, 1)
				}
			}
			io.WriteString(r, line)
			return nil
		})
		if err := scanner.Err(); err != nil {
			return err
		}
	}
	state := seq.finalState()
	if state.exitCode != 0 {
		panic("we aren't using sequencer for any errors")
	}
	// Note that all goroutines have stopped at this point.
	if atomicFixedCount > 0 && path != "-" {
		in.Close()
		// Overwrite the file, if we weren't reading stdin. Report its
		// path too.
		fmt.Println(path)
		if err := ioutil.WriteFile(path, outBuf.Bytes(), 0o666); err != nil {
			return err
		}
	}
	if len(state.brokenURLs) > 0 {
		return fmt.Errorf("found %d broken urls in %q:\n%s", len(state.brokenURLs),
			path, strings.Join(state.brokenURLs, "\n"))
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
