/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package xurls

import (
	"testing"
)

func TestWebUrl(t *testing.T) {
	for _, c := range [...]struct {
		in   string
		want string
	}{
		{"", ""},
		{"foo", ""},
		{"foo.a", ""},
		{"foo.bar", "foo.bar"},
		{"foo.bar/", "foo.bar/"},
		{"1.1.1.1", "1.1.1.1"},
		{"121.1.1.1", "121.1.1.1"},
		{"255.1.1.1", "255.1.1.1"},
		{"300.1.1.1", ""},
		{"1.1.1", ""},
		{"1.1..1", ""},
		{"test.foo.bar", "test.foo.bar"},
		{"test.foo.bar/path", "test.foo.bar/path"},
		{"test.foo.bar/path_(more)", "test.foo.bar/path_(more)"},
		{"http://foo.bar", "http://foo.bar"},
		{" http://foo.bar ", "http://foo.bar"},
		{",http://foo.bar,", "http://foo.bar"},
		{"(http://foo.bar)", "http://foo.bar"},
		{"<http://foo.bar>", "http://foo.bar"},
		{"\"http://foo.bar\"", "http://foo.bar"},
		{"http://foo.bar", "http://foo.bar"},
		{"http://test.foo.bar/", "http://test.foo.bar/"},
		{"http://foo.bar/path", "http://foo.bar/path"},
		{"http://1.1.1.1/path", "http://1.1.1.1/path"},
		{"www.foo.bar", "www.foo.bar"},
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
		{"foo@bar.tld", "foo@bar.tld"},
		{"mailto:foo@bar.tld", "mailto:foo@bar.tld"},
		{"foo@test.bar.tld", "foo@test.bar.tld"},
		{"foo@bar.tld/path", "foo@bar.tld"},
		{"foo+test@bar.tld", "foo+test@bar.tld"},
		{"foo+._%-@bar.tld", "foo+._%-@bar.tld"},
	} {
		got := EmailAddr.FindString(c.in)
		if got != c.want {
			t.Errorf(`EmailAddr.FindString("%s") got "%s", want "%s"`, c.in, got, c.want)
		}
	}
}
