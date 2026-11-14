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
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"mvdan.cc/xurls/v2"
)

var (
	matching    = flag.String("m", "", "")
	relaxed     = flag.Bool("r", false, "")
	fix         boolString
	versionFlag = flag.Bool("version", false, "")
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
Usage: xurls [flags] [files]

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

	seq := newSequencer(32, out) // a maximum of 32 parallel requests

	userAgent := fmt.Sprintf("mvdan.cc/xurls %s", readVersion())

	// Without -fix, buffer the printed urls.
	var input io.Reader = in
	var stdout *bufio.Writer
	if fix == "" {
		stdout = bufio.NewWriter(os.Stdout)
		// Flush before reads which may block, such as from a pipe.
		if info, err := in.Stat(); err != nil || !info.Mode().IsRegular() {
			input = flushReader{in, stdout}
		}
	}
	br := bufio.NewReader(input)

	// Doesn't need to be part of reporterState as order doesn't matter.
	var fixedCount atomic.Uint32

	var readErr error
	for readErr == nil {
		// Keep any line ending as-is, so that -fix leaves it unchanged.
		var line string
		line, readErr = br.ReadString('\n')
		if line == "" {
			continue // only along with io.EOF or another error
		}
		matches := re.FindAllStringIndex(line, -1)
		if fix == "" {
			for _, pair := range matches {
				match := line[pair[0]:pair[1]]
				fmt.Fprintln(stdout, match)
			}
			continue
		}
		seq.Add(int64(len(matches)), func(r *reporter) {
			var fixedLine strings.Builder
			last := 0 // end of the last replaced url in line
			for _, pair := range matches {
				match := line[pair[0]:pair[1]]
				origURL, err := url.Parse(match)
				if err != nil {
					r.appendBroken(match, err.Error())
					continue
				}
				fixed := match
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
							fixed = req.URL.String()
							if req.URL.Fragment == "" && origURL.Fragment != "" {
								// Inherit the fragment as written, without escaping it.
								fixed += match[strings.IndexByte(match, '#'):]
							}
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
					fixedLine.WriteString(line[last:pair[0]])
					fixedLine.WriteString(fixed)
					last = pair[1]
					fixedCount.Add(1)
				}
			}
			if last > 0 {
				fixedLine.WriteString(line[last:])
				line = fixedLine.String()
			}
			io.WriteString(r, line) // add the fixed line to outBuf
		})
	}
	if stdout != nil {
		if err := stdout.Flush(); err != nil {
			return err
		}
	}
	state := seq.finalState()
	var brokenErr error
	if len(state.brokenURLs) > 0 {
		var s strings.Builder
		fmt.Fprintf(&s, "found %d broken urls in %q:", len(state.brokenURLs), path)
		for _, broken := range state.brokenURLs {
			fmt.Fprintf(&s, "\n  * %s - %s", broken.url, broken.reason)
		}
		brokenErr = errors.New(s.String())
	}
	if readErr != io.EOF {
		// Don't overwrite the file with partial output.
		return errors.Join(readErr, brokenErr)
	}
	// Note that all goroutines have stopped at this point.
	if fixedCount.Load() > 0 && path != "-" {
		in.Close()
		// Overwrite the file, if we weren't reading stdin. Report its
		// path too.
		fmt.Println(path)
		if err := os.WriteFile(path, outBuf.Bytes(), 0o666); err != nil {
			return err
		}
	}
	return brokenErr
}

// flushReader flushes a writer before each read.
type flushReader struct {
	r io.Reader
	w *bufio.Writer
}

func (fr flushReader) Read(p []byte) (int, error) {
	if err := fr.w.Flush(); err != nil {
		return 0, err
	}
	return fr.r.Read(p)
}

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(readVersion())
		return
	}
	if *relaxed && *matching != "" {
		fmt.Fprintln(os.Stderr, "-r and -m at the same time don't make much sense")
		os.Exit(1)
	}
	switch fix {
	case "": // disabled by default
	case "false": // disabled via -fix=false; normalize
		fix = ""
	case "auto", "all": // enabled via -fix=auto, -fix=all, etc
	case "true": // enabled via -fix; normalize
		fix = "auto"
	default:
		flag.Usage()
		os.Exit(2)
	}
	var re *regexp.Regexp
	if *relaxed {
		re = xurls.Relaxed()
	} else if *matching != "" {
		var err error
		if re, err = xurls.StrictMatchingScheme(*matching); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	} else {
		re = xurls.Strict()
	}
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}
	failed := false
	for _, path := range args {
		if err := scanPath(re, path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
}

func readVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(unknown)"
	}
	mod := &info.Main
	if mod.Replace != nil {
		mod = mod.Replace
	}
	return mod.Version
}
