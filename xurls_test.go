// Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package xurls

import (
	"regexp"
	"testing"
)

type testCase struct {
	in   string
	want interface{}
}

func doTest(t *testing.T, name string, re *regexp.Regexp, cases []testCase) {
	for _, c := range cases {
		got := re.FindString(c.in)
		want, _ := c.want.(string)
		if got != want {
			t.Errorf(`%s.FindString("%s") got "%s", want "%s"`, name, c.in, got, want)
		}
	}
}

var constantTestCases = []testCase{
	{``, nil},
	{` `, nil},
	{`:`, nil},
	{`::`, nil},
	{`:::`, nil},
	{`::::`, nil},
	{`.`, nil},
	{`..`, nil},
	{`...`, nil},
	{`1.1`, nil},
	{`.1.`, nil},
	{`1.1.1`, nil},
	{`1:1`, nil},
	{`:1:`, nil},
	{`1:1:1`, nil},
	{`://`, nil},
	{`foo`, nil},
	{`foo:`, nil},
	{`mailto:`, nil},
	{`foo://`, nil},
	{`http://`, nil},
	{`http:// foo`, nil},
	{`http:// foo`, nil},
	{`:foo`, nil},
	{`://foo`, nil},
	{`foorandom:bar`, nil},
	{`foo.randombar`, nil},
	{`zzz.`, nil},
	{`.zzz`, nil},
	{`zzz.zzz`, nil},
	{`/some/path`, nil},
	{`rel/path`, nil},
	{`localhost`, nil},
	{`com`, nil},
	{`.com`, nil},
	{`com.`, nil},
	{`http`, nil},

	{`http://foo`, `http://foo`},
	{`http://FOO`, `http://FOO`},
	{`http://FAÀ`, `http://FAÀ`},
	{`https://localhost`, `https://localhost`},
	{`git+https://localhost`, `git+https://localhost`},
	{`foo.bar://localhost`, `foo.bar://localhost`},
	{`foo-bar://localhost`, `foo-bar://localhost`},
	{`mailto:foo`, `mailto:foo`},
	{`MAILTO:foo`, `MAILTO:foo`},
	{`sms:123`, `sms:123`},
	{`xmpp:foo@bar`, `xmpp:foo@bar`},
	{`bitcoin:Addr23?amount=1&message=foo`, `bitcoin:Addr23?amount=1&message=foo`},
	{`http://foo.com`, `http://foo.com`},
	{`http://foo.co.uk`, `http://foo.co.uk`},
	{`http://foo.random`, `http://foo.random`},
	{` http://foo.com/bar `, `http://foo.com/bar`},
	{` http://foo.com/bar more`, `http://foo.com/bar`},
	{`<http://foo.com/bar>`, `http://foo.com/bar`},
	{`<http://foo.com/bar>more`, `http://foo.com/bar`},
	{`.http://foo.com/bar.`, `http://foo.com/bar`},
	{`.http://foo.com/bar.more`, `http://foo.com/bar.more`},
	{`,http://foo.com/bar,`, `http://foo.com/bar`},
	{`,http://foo.com/bar,more`, `http://foo.com/bar,more`},
	{`(http://foo.com/bar)`, `http://foo.com/bar`},
	{`"http://foo.com/bar'`, `http://foo.com/bar`},
	{`"http://foo.com/bar'more`, `http://foo.com/bar'more`},
	{`"http://foo.com/bar"`, `http://foo.com/bar`},
	{`http://a.b/a0/-+_&~*%@|=#.,:;'?![]()a`, `http://a.b/a0/-+_&~*%@|=#.,:;'?![]()a`},
	{`http://foo.bar/path/`, `http://foo.bar/path/`},
	{`http://foo.bar/path-`, `http://foo.bar/path-`},
	{`http://foo.bar/path+`, `http://foo.bar/path+`},
	{`http://foo.bar/path_`, `http://foo.bar/path_`},
	{`http://foo.bar/path&`, `http://foo.bar/path&`},
	{`http://foo.bar/path~`, `http://foo.bar/path~`},
	{`http://foo.bar/path*`, `http://foo.bar/path*`},
	{`http://foo.bar/path%`, `http://foo.bar/path%`},
	{`http://foo.bar/path$`, `http://foo.bar/path$`},
	{`http://foo.bar/path€`, `http://foo.bar/path€`},
	{`http://foo.bar/path@`, `http://foo.bar/path`},
	{`http://foo.bar/path|`, `http://foo.bar/path`},
	{`http://foo.bar/path=`, `http://foo.bar/path`},
	{`http://foo.bar/path#`, `http://foo.bar/path`},
	{`http://foo.bar/path.`, `http://foo.bar/path`},
	{`http://foo.bar/path,`, `http://foo.bar/path`},
	{`http://foo.bar/path:`, `http://foo.bar/path`},
	{`http://foo.bar/path;`, `http://foo.bar/path`},
	{`http://foo.bar/path'`, `http://foo.bar/path`},
	{`http://foo.bar/path?`, `http://foo.bar/path`},
	{`http://foo.bar/path!`, `http://foo.bar/path`},
	{`http://foo.bar/path´`, `http://foo.bar/path`},
	{`http://foo.com/path_(more)`, `http://foo.com/path_(more)`},
	{`(http://foo.com/path_(more))`, `http://foo.com/path_(more)`},
	{`http://foo.com/path_(even)-(more)`, `http://foo.com/path_(even)-(more)`},
	{`http://foo.com/path_(even)(more)`, `http://foo.com/path_(even)(more)`},
	{`http://foo.com/path_(even_(nested))`, `http://foo.com/path_(even_(nested))`},
	{`(http://foo.com/path_(even_(nested)))`, `http://foo.com/path_(even_(nested))`},
	{`http://foo.com/path_[more]`, `http://foo.com/path_[more]`},
	{`[http://foo.com/path_[more]]`, `http://foo.com/path_[more]`},
	{`http://foo.com/path_[even]-[more]`, `http://foo.com/path_[even]-[more]`},
	{`http://foo.com/path_[even][more]`, `http://foo.com/path_[even][more]`},
	{`http://foo.com/path_[even_[nested]]`, `http://foo.com/path_[even_[nested]]`},
	{`[http://foo.com/path_[even_[nested]]]`, `http://foo.com/path_[even_[nested]]`},
	{`http://foo.com/path#fragment`, `http://foo.com/path#fragment`},
	{`http://test.foo.com/`, `http://test.foo.com/`},
	{`http://foo.com/path`, `http://foo.com/path`},
	{`http://foo.com:8080/path`, `http://foo.com:8080/path`},
	{`http://1.1.1.1/path`, `http://1.1.1.1/path`},
	{`http://1080::8:800:200c:417a/path`, `http://1080::8:800:200c:417a/path`},
	{`http://中国.中国/foo中国`, `http://中国.中国/foo中国`},
	{`http://xn-foo.xn--p1acf/path`, `http://xn-foo.xn--p1acf/path`},
	{`http://✪foo.bar/pa✪th`, `http://✪foo.bar/pa✪th`},
	{`✪http://✪foo.bar/pa✪th✪`, `http://✪foo.bar/pa✪th`},
	{`what is http://foo.com?`, `http://foo.com`},
	{`go visit http://foo.com/path.`, `http://foo.com/path`},
	{`go visit http://foo.com/path...`, `http://foo.com/path`},
	{`what is http://foo.com/path?`, `http://foo.com/path`},
	{`the http://foo.com!`, `http://foo.com`},
	{`https://test.foo.bar/path?a=b`, `https://test.foo.bar/path?a=b`},
	{`ftp://user@foo.bar`, `ftp://user@foo.bar`},
	{`http://foo.com/@"style="color:red"onmouseover=func()`, `http://foo.com/`},
}

func TestRegexes(t *testing.T) {
	doTest(t, "Relaxed", Relaxed, constantTestCases)
	doTest(t, "Strict", Strict, constantTestCases)
	doTest(t, "Relaxed", Relaxed, []testCase{
		{`foo.a`, nil},
		{`foo.com`, `foo.com`},
		{`foo.com bar.com`, `foo.com`},
		{`foo.com-foo`, `foo.com`},
		{`foo.company`, `foo.company`},
		{`foo.comrandom`, nil},
		{`foo.onion`, `foo.onion`},
		{`foo.i2p`, `foo.i2p`},
		{`中国.中国`, `中国.中国`},
		{`中国.中国/foo中国`, `中国.中国/foo中国`},
		{`foo.com/`, `foo.com/`},
		{`1.1.1.1`, `1.1.1.1`},
		{`10.50.23.250`, `10.50.23.250`},
		{`121.1.1.1`, `121.1.1.1`},
		{`255.1.1.1`, `255.1.1.1`},
		{`300.1.1.1`, nil},
		{`1.1.1.300`, nil},
		{`1080:0:0:0:8:800:200C:4171`, `1080:0:0:0:8:800:200C:4171`},
		{`3ffe:2a00:100:7031::1`, `3ffe:2a00:100:7031::1`},
		{`1080::8:800:200c:417a`, `1080::8:800:200c:417a`},
		{`foo.com:8080`, `foo.com:8080`},
		{`foo.com:8080/path`, `foo.com:8080/path`},
		{`test.foo.com`, `test.foo.com`},
		{`test.foo.com/path`, `test.foo.com/path`},
		{`test.foo.com/path/more/`, `test.foo.com/path/more/`},
		{`TEST.FOO.COM/PATH`, `TEST.FOO.COM/PATH`},
		{`TEST.FÓO.COM/PÁTH`, `TEST.FÓO.COM/PÁTH`},
		{`foo.com/path_(more)`, `foo.com/path_(more)`},
		{`foo.com/path_(even)_(more)`, `foo.com/path_(even)_(more)`},
		{`foo.com/path_(more)/more`, `foo.com/path_(more)/more`},
		{`foo.com/path_(more)/end)`, `foo.com/path_(more)/end`},
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
		{`"foo.com/bar'`, `foo.com/bar`},
		{`"foo.com/bar'more`, `foo.com/bar'more`},
		{`"foo.com/bar"`, `foo.com/bar`},
		{`what is foo.com?`, `foo.com`},
		{`the foo.com!`, `foo.com`},

		{`foo@bar`, nil},
		{`foo@bar.a`, nil},
		{`foo@bar.com`, `foo@bar.com`},
		{`foo@bar.com bar@bar.com`, `foo@bar.com`},
		{`foo@bar.onion`, `foo@bar.onion`},
		{`foo@中国.中国`, `foo@中国.中国`},
		{`foo@test.bar.com`, `foo@test.bar.com`},
		{`FOO@TEST.BAR.COM`, `FOO@TEST.BAR.COM`},
		{`foo@bar.com/path`, `foo@bar.com`},
		{`foo+test@bar.com`, `foo+test@bar.com`},
		{`foo+._%-@bar.com`, `foo+._%-@bar.com`},
	})
	doTest(t, "Strict", Strict, []testCase{
		{`http:// foo.com`, nil},
		{`foo.a`, nil},
		{`foo.com`, nil},
		{`foo.com/`, nil},
		{`1.1.1.1`, nil},
		{`3ffe:2a00:100:7031::1`, nil},
		{`test.foo.com:8080/path`, nil},
		{`foo@bar.com`, nil},
	})
}

func TestStrictMatchingError(t *testing.T) {
	for _, c := range []struct {
		exp     string
		wantErr bool
	}{
		{`http://`, false},
		{`https?://`, false},
		{`http://|mailto:`, false},
		{`http://(`, true},
	} {
		_, err := StrictMatching(c.exp)
		if c.wantErr && err == nil {
			t.Errorf(`StrictMatching("%s") did not error as expected`, c.exp)
		} else if !c.wantErr && err != nil {
			t.Errorf(`StrictMatching("%s") unexpectedly errored`, c.exp)
		}
	}
}

func TestStrictMatching(t *testing.T) {
	strictMatching, _ := StrictMatching("http://|ftps?://|mailto:")
	doTest(t, "StrictMatching", strictMatching, []testCase{
		{`foo.com`, nil},
		{`foo@bar.com`, nil},
		{`http://foo`, `http://foo`},
		{`Http://foo`, `Http://foo`},
		{`https://foo`, nil},
		{`ftp://foo`, `ftp://foo`},
		{`ftps://foo`, `ftps://foo`},
		{`mailto:foo`, `mailto:foo`},
		{`MAILTO:foo`, `MAILTO:foo`},
		{`sms:123`, nil},
	})
}

func bench(b *testing.B, re *regexp.Regexp, str string) {
	for i := 0; i < b.N; i++ {
		re.FindAllString(str, -1)
	}
}

func BenchmarkStrictEmpty(b *testing.B) {
	bench(b, Strict, "foo")
}

func BenchmarkStrictSingle(b *testing.B) {
	bench(b, Strict, "http://foo.foo foo.com")
}

func BenchmarkStrictMany(b *testing.B) {
	bench(b, Strict, ` foo bar http://foo.foo
	foo.com bitcoin:address ftp://
	xmpp:foo@bar.com`)
}

func BenchmarkRelaxedEmpty(b *testing.B) {
	bench(b, Relaxed, "foo")
}

func BenchmarkRelaxedSingle(b *testing.B) {
	bench(b, Relaxed, "http://foo.foo foo.com")
}

func BenchmarkRelaxedMany(b *testing.B) {
	bench(b, Relaxed, ` foo bar http://foo.foo
	foo.com bitcoin:address ftp://
	xmpp:foo@bar.com`)
}
