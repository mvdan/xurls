/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package xurls

import (
	"testing"
)

func TestReverseJoin(t *testing.T) {
	for _, c := range [...]struct {
		in    []string
		inSep string
		want  string
	}{
		{[]string{"a", "b", "c"}, ".", "c.b.a"},
		{[]string{"1", "22", "333"}, ",", "333,22,1"},
		{[]string{"xxx", "yyy", "zzz"}, "__", "zzz__yyy__xxx"},
	} {
		got := reverseJoin(c.in, c.inSep)
		if got != c.want {
			t.Errorf(`reverseJoin(%v) got "%s", want "%s"`, c.in, got, c.want)
		}
	}
}

func TestWebUrl(t *testing.T) {
	for _, c := range [...]struct {
		in   string
		want string
	}{
		{"", ""},
		{"foo", ""},
		{"foo.a", ""},
		{"foo.random", ""},
		{"foo.com", "foo.com"},
		{"foo.com/", "foo.com/"},
		{"1.1.1.1", "1.1.1.1"},
		{"121.1.1.1", "121.1.1.1"},
		{"255.1.1.1", "255.1.1.1"},
		{"300.1.1.1", ""},
		{"1.1.1", ""},
		{"1.1..1", ""},
		{"test.foo.com", "test.foo.com"},
		{"test.foo.com/path", "test.foo.com/path"},
		{"test.foo.com/path_(more)", "test.foo.com/path_(more)"},
		{"http://foo.com", "http://foo.com"},
		{" http://foo.com ", "http://foo.com"},
		{",http://foo.com,", "http://foo.com"},
		{"(http://foo.com)", "http://foo.com"},
		{"<http://foo.com>", "http://foo.com"},
		{"\"http://foo.com\"", "http://foo.com"},
		{"http://foo.com", "http://foo.com"},
		{"http://test.foo.com/", "http://test.foo.com/"},
		{"http://foo.com/path", "http://foo.com/path"},
		{"http://1.1.1.1/path", "http://1.1.1.1/path"},
		{"www.foo.com", "www.foo.com"},
		{" foo.com/bar ", "foo.com/bar"},
		{",foo.com/bar,", "foo.com/bar,"},
		{"(foo.com/bar)", "foo.com/bar)"},
		{"<foo.com/bar>", "foo.com/bar"},
		{"\"foo.com/bar\"", "foo.com/bar"},
	} {
		got := WebUrl.FindString(c.in)
		if got != c.want {
			t.Errorf(`WebUrl.FindString("%s") got "%s", want "%s"`, c.in, got, c.want)
		}
	}
}

func TestEmailAddr(t *testing.T) {
	for _, c := range [...]struct {
		in   string
		want string
	}{
		{"", ""},
		{"foo", ""},
		{"foo@bar", ""},
		{"foo@bar.a", ""},
		{"foo@bar.com", "foo@bar.com"},
		{"mailto:foo@bar.com", "mailto:foo@bar.com"},
		{"foo@test.bar.com", "foo@test.bar.com"},
		{"foo@bar.com/path", "foo@bar.com"},
		{"foo+test@bar.com", "foo+test@bar.com"},
		{"foo+._%-@bar.com", "foo+._%-@bar.com"},
	} {
		got := EmailAddr.FindString(c.in)
		if got != c.want {
			t.Errorf(`EmailAddr.FindString("%s") got "%s", want "%s"`, c.in, got, c.want)
		}
	}
}
