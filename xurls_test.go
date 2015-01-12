/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package xurls

import (
	"testing"
)

func TestFindString(t *testing.T) {
	for _, c := range []struct {
		in   string
		want string
	}{
		{"", ""},
		{"foo", ""},
		{"foo.bar", ""},
		{"test.foo.bar", ""},
		{"test.foo.bar/path", ""},
		{"http://foo.bar", "http://foo.bar"},
		{" foo.com ", "foo.com"},
		{",foo.com,", "foo.com"},
		{"(foo.com)", "foo.com"},
		{"<foo.com>", "foo.com"},
		{"\"foo.com\"", "foo.com"},
		{"http://foo.bar/", "http://foo.bar/"},
		{"http://foo.bar/path", "http://foo.bar/path"},
		{"mailto:foo@bar", "mailto:foo@bar"},
		{" mailto:foo@bar ", "mailto:foo@bar"},
		{",mailto:foo@bar,", "mailto:foo@bar"},
		{"(mailto:foo@bar)", "mailto:foo@bar"},
		{"<mailto:foo@bar>", "mailto:foo@bar"},
		{"\"mailto:foo@bar\"", "mailto:foo@bar"},
		{"www.foo.bar", "www.foo.bar"},
		{" www.foo.bar ", "www.foo.bar"},
		{",www.foo.bar,", "www.foo.bar"},
		{"(www.foo.bar)", "www.foo.bar"},
		{"<www.foo.bar>", "www.foo.bar"},
		{"\"www.foo.bar\"", "www.foo.bar"},
		{"foo.com", "foo.com"},
		{"foo.com/", "foo.com/"},
		{"foo.org/bar", "foo.org/bar"},
		{" foo.com ", "foo.com"},
		{",foo.com,", "foo.com"},
		{"(foo.com)", "foo.com"},
		{"<foo.com>", "foo.com"},
		{"\"foo.com\"", "foo.com"},
	} {
		got := FindString(c.in)
		if got != c.want {
			t.Errorf("FindString(\"%s\") got \"%s\", want \"%s\"", c.in, got, c.want)
		}
	}
}

func stringSliceEqual(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFindAllString(t *testing.T) {
	for _, c := range []struct {
		in   string
		want []string
	}{
		{"", nil},
		{"http://foo.bar", []string{"http://foo.bar"}},
		{" http://foo.bar www.foo.bar ", []string{"http://foo.bar", "www.foo.bar"}},
	} {
		got := FindAllString(c.in)
		if !stringSliceEqual(got, c.want) {
			t.Errorf("urlsFromString(\"%s\") got %q, want %q", c.in, got, c.want)
		}
	}
}
