/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package xurls

import (
	"testing"
	"reflect"
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
		got := Regexp.FindString(c.in)
		if got != c.want {
			t.Errorf(`Regexp.FindString("%s") got "%s", want "%s"`, c.in, got, c.want)
		}
	}
}

func TestFindAllString(t *testing.T) {
	for _, c := range []struct {
		in   string
		inN  int
		want []string
	}{
		{"", -1, nil},
		{"http://foo.bar", 0, nil},
		{"http://foo.bar", -1, []string{"http://foo.bar"}},
		{" http://foo.bar www.foo.bar ", -1, []string{"http://foo.bar", "www.foo.bar"}},
	} {
		got := Regexp.FindAllString(c.in, c.inN)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf(`Regexp.FindAllString("%s") got "%q", want "%q"`, c.in, got, c.want)
		}
	}
}
