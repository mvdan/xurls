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
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/mod/module"

	"mvdan.cc/xurls/v2"
)

var (
	matching = flag.String("m", "", "")
	relaxed  = flag.Bool("r", false, "")
	fix      boolString
	version  = flag.Bool("version", false, "")
)

type boolString string

func (s *boolString) Set(val string) error {
	*s = boolString(val)
	return nil
}
func (s *boolString) Get() any       { return string(*s) }
func (s *boolString) String() string { return string(*s) }
func (*boolString) IsBoolFlag() bool { return true }

func init() {
	flag.Var(&fix, "fix", "")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `
Usage: xurls [-h] [files]

xurls extracts urls from text using regular expressions.
If no files are given, it reads from standard input.

   -m <regexp>   only match urls whose scheme matches a regexp
                    example: 'https?://|mailto:'
   -r            also match urls without a scheme (relaxed)
   -version      print version and exit

When the -fix or -fix=auto flag is used, xurls instead attempts to replace
any urls which result in a permanent redirect (301 or 308).
It also fails if any urls fail to load, so that they may be removed or replaced.
To replace urls which result in temporary redirect as well, use -fix=all.
`[1:])
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
		if fix != "" {
			outBuf = new(bytes.Buffer)
			out = outBuf
		}
		defer in.Close()
	}

	// A maximum of 32 parallel requests.
	maxWeight := int64(32)
	seq := newSequencer(maxWeight, out, os.Stderr)

	userAgent := fmt.Sprintf("mvdan.cc/xurls %s", readVersion())
	scanner := bufio.NewScanner(in)

	// Doesn't need to be part of reporterState as order doesn't matter.
	var atomicFixedCount uint32

	for scanner.Scan() {
		line := scanner.Text() + "\n"
		matches := re.FindAllStringIndex(line, -1)
		if fix == "" {
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
			for _, pair := range matches {
				// The indexes are based on the original line.
				pair[0] += offsetWithinLine
				pair[1] += offsetWithinLine
				match := line[pair[0]:pair[1]]
				origURL, err := url.Parse(match)
				if err != nil {
					r.appendBroken(match, err.Error())
					continue
				}
				fixed := origURL.String()
				switch origURL.Scheme {
				case "http", "https":
					// See if the URL redirects somewhere.
					client := &http.Client{
						Timeout: 10 * time.Second,
						CheckRedirect: func(req *http.Request, via []*http.Request) error {
							if len(via) >= 10 {
								return errors.New("stopped after 10 redirects")
							}
							switch req.Response.StatusCode {
							case http.StatusMovedPermanently, http.StatusPermanentRedirect:
								// "auto" and "all" fix permanent redirects.
							case http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
								// Only "all" fixes temporary redirects.
								if fix != "all" {
									return http.ErrUseLastResponse
								}
							default:
								// Any other redirects are ignored.
								return http.ErrUseLastResponse
							}
							// Inherit the fragment if empty.
							if req.URL.Fragment == "" {
								req.URL.Fragment = origURL.Fragment
							}
							fixed = req.URL.String()
							return nil
						},
					}
					method := http.MethodHead
				retry:
					req, err := http.NewRequest(method, fixed, nil)
					if err != nil {
						r.appendBroken(match, err.Error())
						continue
					}
					req.Header.Set("User-Agent", userAgent)
					resp, err := client.Do(req)
					if err != nil {
						r.appendBroken(match, err.Error())
						continue
					}
					if code := resp.StatusCode; code >= 400 {
						if code == http.StatusMethodNotAllowed {
							method = http.MethodGet
							resp.Body.Close()
							goto retry
						}
						r.appendBroken(match, fmt.Sprintf("%d %s", code, http.StatusText(code)))
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
			io.WriteString(r, line) // add the fixed line to outBuf
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
		var s strings.Builder
		fmt.Fprintf(&s, "found %d broken urls in %q:\n", len(state.brokenURLs), path)
		for _, broken := range state.brokenURLs {
			fmt.Fprintf(&s, "  * %s - %s\n", broken.url, broken.reason)
		}
		return errors.New(s.String())
	}
	return nil
}

func main() { os.Exit(main1()) }

func main1() int {
	flag.Parse()
	if *version {
		fmt.Println(readVersion())
		return 0
	}
	if *relaxed && *matching != "" {
		fmt.Fprintln(os.Stderr, "-r and -m at the same time don't make much sense")
		return 1
	}
	switch fix {
	case "": // disabled by default
	case "false": // disabled via -fix=false; normalize
		fix = ""
	case "auto", "all": // enabled via -fix=auto, -fix=all, etc
	case "true": // enabled via -fix; normalize
		fix = "auto"
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

// Borrowed from https://github.com/burrowers/garble.

func readVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	mod := &info.Main
	if mod.Replace != nil {
		mod = mod.Replace
	}

	// Until https://github.com/golang/go/issues/50603 is implemented,
	// manually construct something like a pseudo-version.
	// TODO: remove when this code is dead, hopefully in Go 1.20.
	if mod.Version == "(devel)" {
		var vcsTime time.Time
		var vcsRevision string
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.time":
				// If the format is invalid, we'll print a zero timestamp.
				vcsTime, _ = time.Parse(time.RFC3339Nano, setting.Value)
			case "vcs.revision":
				vcsRevision = setting.Value
				if len(vcsRevision) > 12 {
					vcsRevision = vcsRevision[:12]
				}
			}
		}
		if vcsRevision != "" {
			mod.Version = module.PseudoVersion("", "", vcsTime, vcsRevision)
		}
	}
	return mod.Version
}
