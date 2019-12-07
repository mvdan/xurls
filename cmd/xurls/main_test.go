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

func TestScripts(t *testing.T) {
	t.Parallel()
	testscript.Run(t, testscript.Params{
		Dir: filepath.Join("testdata", "scripts"),
		Setup: func(env *testscript.Env) error {
			mux := http.NewServeMux()
			mux.HandleFunc("/plain", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "plaintext")
			})
			mux.HandleFunc("/redir-1", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/plain", http.StatusMovedPermanently)
			})
			mux.HandleFunc("/redir-2", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/redir-1", http.StatusMovedPermanently)
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
					err := ioutil.WriteFile(ts.MkAbs(arg), []byte(data), 0666)
					ts.Check(err)
				}
			},
		},
	})
}
