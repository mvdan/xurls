// Copyright (c) 2019, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"xurls": main1,
	}))
}

func TestScript(t *testing.T) {
	t.Parallel()
	testscript.Run(t, testscript.Params{
		Dir:                 filepath.Join("testdata", "script"),
		RequireExplicitExec: true,
		Setup: func(env *testscript.Env) error {
			mux := http.NewServeMux()
			handle := func(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
				mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
					if r.Method != method {
						t.Errorf("expected all requests to be %q, got %q", method, r.Method)
					}
					handler(w, r)
				})
			}
			handle("HEAD", "/plain-head", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			handle("HEAD", "/redir-1", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/plain-head", http.StatusMovedPermanently)
			})
			handle("HEAD", "/redir-2", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/redir-1", http.StatusMovedPermanently)
			})

			handle("HEAD", "/redir-longer", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/redir-longtarget", http.StatusMovedPermanently)
			})
			handle("HEAD", "/redir-longtarget", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			handle("HEAD", "/redir-fragment", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/plain-head#bar", http.StatusMovedPermanently)
			})

			handle("HEAD", "/redir-301", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/plain-head", 301)
			})
			handle("HEAD", "/redir-302", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/plain-head", 302)
			})
			handle("HEAD", "/redir-303", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/plain-head", 303)
			})
			handle("HEAD", "/redir-307", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/plain-head", 307)
			})
			handle("HEAD", "/redir-308", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/plain-head", 308)
			})

			handle("HEAD", "/404", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "", 404)
			})
			handle("HEAD", "/500", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "", 500)
			})

			handle("GET", "/plain-get", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "plaintext")
			})
			mux.HandleFunc("/get-only", func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					http.Redirect(w, r, "/plain-get", 301)
				} else {
					http.Error(w, "", 405)
				}
			})

			ln, err := net.Listen("tcp", ":0")
			if err != nil {
				return err
			}
			server := &http.Server{Handler: mux}
			go server.Serve(ln)
			env.Vars = append(env.Vars, "SERVER=http://"+ln.Addr().String())
			env.Defer(func() {
				if err := server.Shutdown(context.TODO()); err != nil {
					t.Fatal(err)
				}
			})
			return nil
		},
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"expand": func(ts *testscript.TestScript, neg bool, args []string) {
				if neg {
					ts.Fatalf("unsupported: ! expand")
				}
				if len(args) == 0 {
					ts.Fatalf("usage: expand file...")
				}
				for _, arg := range args {
					data := ts.ReadFile(arg)
					data = os.Expand(data, ts.Getenv)
					err := ioutil.WriteFile(ts.MkAbs(arg), []byte(data), 0o666)
					ts.Check(err)
				}
			},
		},
	})
}
