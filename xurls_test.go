/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package xurls

import "testing"

func TestWebURL(t *testing.T) {
	for _, c := range [...]struct {
		in   string
		want interface{}
	}{
		{``, nil},
		{` `, nil},
		{`:`, nil},
		{`::`, nil},
		{`:::`, nil},
		{`::::`, nil},
		{`.`, nil},
		{`..`, nil},
		{`...`, nil},
		{`://`, nil},
		{`foo`, nil},
		{`foo:`, nil},
		{`foo://`, nil},
		{`:foo`, nil},
		{`://foo`, nil},
		{`foo:bar`, nil},
		{`/some/path`, nil},

		// Web links
		{`foo.a`, nil},
		{`foo.com`, `foo.com`},
		{`foo.com bar.com`, `foo.com`},
		{`foo.com-foo`, `foo.com`},
		{`foo.onion`, `foo.onion`},
		{`foo.i2p`, `foo.i2p`},
		{`中国.中国`, `中国.中国`},
		{`中国.中国/foo中国`, `中国.中国/foo中国`},
		{`foo.com/`, `foo.com/`},
		{`1.1.1.1`, `1.1.1.1`},
		{`121.1.1.1`, `121.1.1.1`},
		{`255.1.1.1`, `255.1.1.1`},
		{`300.1.1.1`, nil},
		{`1.1.1`, nil},
		{`1.1..1`, nil},
		{`1080:0:0:0:8:800:200C:4171`, `1080:0:0:0:8:800:200C:4171`},
		{`3ffe:2a00:100:7031::1`, `3ffe:2a00:100:7031::1`},
		{`1080::8:800:200c:417a`, `1080::8:800:200c:417a`},
		{`1:1`, nil},
		{`:2:`, nil},
		{`1:2:3`, nil},
		{`foo.com:8080`, `foo.com:8080`},
		{`foo.com:8080/path`, `foo.com:8080/path`},
		{`test.foo.com`, `test.foo.com`},
		{`test.foo.com/path`, `test.foo.com/path`},
		{`test.foo.com/path/more/`, `test.foo.com/path/more/`},
		{`TEST.FOO.COM/PATH`, `TEST.FOO.COM/PATH`},
		{`foo.com/a.,:;-+_()?@&=$~!*%'"a`, `foo.com/a.,:;-+_()?@&=$~!*%'"a`},
		{`foo.com/path_(more)`, `foo.com/path_(more)`},
		{`foo.com/path_(even)_(more)`, `foo.com/path_(even)_(more)`},
		{`foo.com/path_(more)/more`, `foo.com/path_(more)/more`},
		{`foo.com/path_(more)/end)`, `foo.com/path_(more)/end)`},
		{`www.foo.com`, `www.foo.com`},
		{` foo.com/bar `, `foo.com/bar`},
		{` foo.com/bar more`, `foo.com/bar`},
		{`<foo.com/bar>`, `foo.com/bar`},
		{`<foo.com/bar>more`, `foo.com/bar`},
		{`,foo.com/bar.`, `foo.com/bar`},
		{`,foo.com/bar.more`, `foo.com/bar.more`},
		{`,foo.com/bar,`, `foo.com/bar`},
		{`,foo.com/bar,more`, `foo.com/bar,more`},
		{`(foo.com/bar)`, `foo.com/bar`},
		{`(foo.com/bar)more`, `foo.com/bar)more`},
		{`"foo.com/bar'`, `foo.com/bar`},
		{`"foo.com/bar'more`, `foo.com/bar'more`},
		{`"foo.com/bar"`, `foo.com/bar`},
		{`"foo.com/bar"more`, `foo.com/bar"more`},
		{`what is foo.com?`, `foo.com`},
		{`the foo.com!`, `foo.com`},

		// Email addresses
		{`foo@bar`, nil},
		{`foo@bar.a`, nil},
		{`foo@bar.com`, `foo@bar.com`},
		{`foo@bar.com bar@bar.com`, `foo@bar.com`},
		{`foo@bar.onion`, `foo@bar.onion`},
		{`foo@中国.中国`, `foo@中国.中国`},
		{`mailto:foo@bar.com`, `foo@bar.com`},
		{`foo@test.bar.com`, `foo@test.bar.com`},
		{`FOO@TEST.BAR.COM`, `FOO@TEST.BAR.COM`},
		{`foo@bar.com/path`, `foo@bar.com`},
		{`foo+test@bar.com`, `foo+test@bar.com`},
		{`foo+._%-@bar.com`, `foo+._%-@bar.com`},

		// Urls with scheme and ://
		{`http://foo.com`, `http://foo.com`},
		{`http://foo.random`, `http://foo.random`},
		{` http://foo.com/bar `, `http://foo.com/bar`},
		{` http://foo.com/bar more`, `http://foo.com/bar`},
		{`<http://foo.com/bar>`, `http://foo.com/bar`},
		{`<http://foo.com/bar>more`, `http://foo.com/bar`},
		{`,http://foo.com/bar.`, `http://foo.com/bar`},
		{`,http://foo.com/bar.more`, `http://foo.com/bar.more`},
		{`,http://foo.com/bar,`, `http://foo.com/bar`},
		{`,http://foo.com/bar,more`, `http://foo.com/bar,more`},
		{`(http://foo.com/bar)`, `http://foo.com/bar`},
		{`(http://foo.com/bar)more`, `http://foo.com/bar)more`},
		{`"http://foo.com/bar'`, `http://foo.com/bar`},
		{`"http://foo.com/bar'more`, `http://foo.com/bar'more`},
		{`"http://foo.com/bar"`, `http://foo.com/bar`},
		{`"http://foo.com/bar"more`, `http://foo.com/bar"more`},
		{`http://a.b/a.,:;-+_()?@&=$~!*%'"a`, `http://a.b/a.,:;-+_()?@&=$~!*%'"a`},
		{`http://foo.com/path_(more)`, `http://foo.com/path_(more)`},
		{`http://test.foo.com/`, `http://test.foo.com/`},
		{`http://foo.com/path`, `http://foo.com/path`},
		{`http://foo.com:8080/path`, `http://foo.com:8080/path`},
		{`http://1.1.1.1/path`, `http://1.1.1.1/path`},
		{`http://1080::8:800:200c:417a/path`, `http://1080::8:800:200c:417a/path`},
		{`http://1080::8:800:200c:417a/path`, `http://1080::8:800:200c:417a/path`},
		{`what is http://foo.com?`, `http://foo.com`},
		{`the http://foo.com!`, `http://foo.com`},
		{`https://test.foo.bar/path?a=b`, `https://test.foo.bar/path?a=b`},
		{`ftp://user@foo.bar`, `ftp://user@foo.bar`},
	} {
		got := All.FindString(c.in)
		var want string
		switch x := c.want.(type) {
		case string:
			want = x
		}
		if got != want {
			t.Errorf(`xurls.All.FindString("%s") got "%s", want "%s"`, c.in, got, want)
		}
	}
}
