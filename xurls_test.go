// Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package xurls

import (
	"fmt"
	"regexp"
	"sync"
	"testing"
)

type testCase struct {
	in   string
	want interface{}
}

func wantStr(in string, want interface{}) string {
	switch x := want.(type) {
	case string:
		return x
	case bool:
		if x {
			return in
		}
	}
	return ""
}

func doTest(t *testing.T, name string, re *regexp.Regexp, cases []testCase) {
	for i, c := range cases {
		t.Run(fmt.Sprintf("%s/%03d", name, i), func(t *testing.T) {
			want := wantStr(c.in, c.want)
			for _, surround := range []string{"", "\n"} {
				in := surround + c.in + surround
				got := re.FindString(in)
				if got != want {
					t.Errorf(`FindString(%q) got %q, want %q`, in, got, want)
				}
			}
		})
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

	{`http://foo`, true},
	{`http://FOO`, true},
	{`http://FAÀ`, true},
	{`https://localhost`, true},
	{`mailto:foo`, true},
	{`MAILTO:foo`, true},
	{`sms:123`, true},
	{`xmpp:foo@bar`, true},
	{`bitcoin:Addr23?amount=1&message=foo`, true},
	{`cid:foo-32x32.v2_fe0f1423.png`, true},
	{`mid:960830.1639@XIson.com`, true},
	{`http://foo.com`, true},
	{`http://foo.co.uk`, true},
	{`http://foo.random`, true},
	{` http://foo.com/bar `, `http://foo.com/bar`},
	{` http://foo.com/bar more`, `http://foo.com/bar`},
	{`<http://foo.com/bar>`, `http://foo.com/bar`},
	{`<http://foo.com/bar>more`, `http://foo.com/bar`},
	{`.http://foo.com/bar.`, `http://foo.com/bar`},
	{`.http://foo.com/bar.more`, `http://foo.com/bar.more`},
	{`,http://foo.com/bar,`, `http://foo.com/bar`},
	{`,http://foo.com/bar,more`, `http://foo.com/bar,more`},
	{`*http://foo.com/bar*`, `http://foo.com/bar`},
	{`*http://foo.com/bar*more`, `http://foo.com/bar*more`},
	{`_http://foo.com/bar_`, `http://foo.com/bar_`},
	{`_http://foo.com/bar_more`, `http://foo.com/bar_more`},
	{`(http://foo.com/bar)`, `http://foo.com/bar`},
	{`(http://foo.com/bar)more`, `http://foo.com/bar`},
	{`[http://foo.com/bar]`, `http://foo.com/bar`},
	{`[http://foo.com/bar]more`, `http://foo.com/bar`},
	{`'http://foo.com/bar'`, `http://foo.com/bar`},
	{`'http://foo.com/bar'more`, `http://foo.com/bar'more`},
	{`"http://foo.com/bar"`, `http://foo.com/bar`},
	{`"http://foo.com/bar"more`, `http://foo.com/bar`},
	{`{"url":"http://foo.com/bar"}`, `http://foo.com/bar`},
	{`{"before":"foo","url":"http://foo.com/bar","after":"bar"}`, `http://foo.com/bar`},
	{`http://a.b/a0/-+_&~*%=#@.,:;'?![]()a`, true},
	{`http://a.b/a0/$€¥`, true},
	{`http://✪foo.bar/pa✪th©more`, true},
	{`http://foo.bar/path/`, true},
	{`http://foo.bar/path-`, true},
	{`http://foo.bar/path+`, true},
	{`http://foo.bar/path&`, true},
	{`http://foo.bar/path~`, true},
	{`http://foo.bar/path%`, true},
	{`http://foo.bar/path=`, true},
	{`http://foo.bar/path#`, true},
	{`http://foo.bar/path.`, `http://foo.bar/path`},
	{`http://foo.bar/path,`, `http://foo.bar/path`},
	{`http://foo.bar/path:`, `http://foo.bar/path`},
	{`http://foo.bar/path;`, `http://foo.bar/path`},
	{`http://foo.bar/path'`, `http://foo.bar/path`},
	{`http://foo.bar/path?`, `http://foo.bar/path`},
	{`http://foo.bar/path!`, `http://foo.bar/path`},
	{`http://foo.bar/path@`, `http://foo.bar/path`},
	{`http://foo.bar/path|`, `http://foo.bar/path`},
	{`http://foo.bar/path|more`, `http://foo.bar/path`},
	{`http://foo.bar/path<`, `http://foo.bar/path`},
	{`http://foo.bar/path<more`, `http://foo.bar/path`},
	{`http://foo.com/path_(more)`, true},
	{`(http://foo.com/path_(more))`, `http://foo.com/path_(more)`},
	{`http://foo.com/path_(even)-(more)`, true},
	{`http://foo.com/path_(even)(more)`, true},
	{`http://foo.com/path_(even_(nested))`, true},
	{`(http://foo.com/path_(even_(nested)))`, `http://foo.com/path_(even_(nested))`},
	{`http://foo.com/path_[more]`, true},
	{`[http://foo.com/path_[more]]`, `http://foo.com/path_[more]`},
	{`http://foo.com/path_[even]-[more]`, true},
	{`http://foo.com/path_[even][more]`, true},
	{`http://foo.com/path_[even_[nested]]`, true},
	{`[http://foo.com/path_[even_[nested]]]`, `http://foo.com/path_[even_[nested]]`},
	{`http://foo.com/path_{more}`, true},
	{`{http://foo.com/path_{more}}`, `http://foo.com/path_{more}`},
	{`http://foo.com/path_{even}-{more}`, true},
	{`http://foo.com/path_{even}{more}`, true},
	{`http://foo.com/path_{even_{nested}}`, true},
	{`{http://foo.com/path_{even_{nested}}}`, `http://foo.com/path_{even_{nested}}`},
	{`http://foo.com/path#fragment`, true},
	{`http://foo.com/emptyfrag#`, true},
	{`http://foo.com/spaced%20path`, true},
	{`http://foo.com/?p=spaced%20param`, true},
	{`http://test.foo.com/`, true},
	{`http://foo.com/path`, true},
	{`http://foo.com:8080/path`, true},
	{`http://1.1.1.1/path`, true},
	{`http://1.1.1.1:8080/path`, true},
	{`http://[1080::8:800:200c:417a]/path`, true},
	{`http://[1080::8:800:200c:417a]:8080/path`, true},

	// scheme://IPv6_addr is not valid per RFC 3987, but is supported anyway (for now).
	{`http://1080::8:800:200c:417a/path`, true},
	{`http://2001.db8:0/path`, true},

	{`http://中国.中国/中国`, true},
	{`http://中国.中国/foo中国`, true},
	{`http://उदाहरण.परीकषा`, true},
	{`http://xn-foo.xn--p1acf/path`, true},
	{`what is http://foo.com?`, `http://foo.com`},
	{`go visit http://foo.com/path.`, `http://foo.com/path`},
	{`go visit http://foo.com/path...`, `http://foo.com/path`},
	{`what is http://foo.com/path?`, `http://foo.com/path`},
	{`the http://foo.com!`, `http://foo.com`},
	{`https://test.foo.bar/path?a=b`, `https://test.foo.bar/path?a=b`},
	{`ftp://user@foo.bar`, true},
	{`http://foo.com/base64-bCBwbGVhcw==`, true},
	{`http://foo.com/–`, true},
	{`http://foo.com/🐼`, true},
	{`https://shmibbles.me/tmp/自殺でも？.png`, true},
	{`randomtexthttp://foo.bar/etc`, "http://foo.bar/etc"},
	{`postgres://user:pass@host.com:5432/path?k=v#f`, true},
	{`postgres://user:pass@host.com:5432/path?k=v#f`, true},
	{`zoommtg://zoom.us/join?confno=1234&pwd=xxx`, true},
	{`zoomus://zoom.us/join?confno=1234&pwd=xxx`, true},
}

func TestRegexes(t *testing.T) {
	doTest(t, "Relaxed", Relaxed(), constantTestCases)
	doTest(t, "Strict", Strict(), constantTestCases)
	doTest(t, "Relaxed2", Relaxed(), []testCase{
		{`foo.a`, nil},
		{`foo.com`, true},
		{`foo.com bar.com`, `foo.com`},
		{`foo.com-foo`, `foo.com`},
		{`foo.company`, true},
		{`foo.comrandom`, nil},
		{`some.guy`, nil},
		{`foo.example`, true},
		{`foo.i2p`, true},
		{`foo.local`, true},
		{`foo.onion`, true},
		{`中国.中国`, true},
		{`中国.中国/foo中国`, true},
		{`test.联通`, true},
		{`test.联通 extra`, `test.联通`},
		{`test.xn--8y0a063a`, true},
		{`test.xn--8y0a063a/foobar`, true},
		{`test.xn-foo`, nil},
		{`test.xn--`, nil},
		{`foo.com/`, true},
		{`1.1.1.1`, true},
		{`10.50.23.250`, true},
		{`121.1.1.1`, true},
		{`255.1.1.1`, true},
		{`300.1.1.1`, nil},
		{`1.1.1.300`, nil},
		{`foo@1.2.3.4`, `1.2.3.4`},

		// https://www.iana.org/assignments/iana-ipv6-special-registry/iana-ipv6-special-registry.xhtml
		{`::1`, true},
		//{`::`, true},
		{`::ffff:0:0`, true},
		{`64:ff9b::`, true},
		{`64:ff9b:1::`, true},
		{`100::`, true},
		{`2001::`, true},
		{`2001:1::1`, true},
		{`2001:1::2`, true},
		{`2001:2::`, true},
		{`2001:3::`, true},
		{`2001:4:112::`, true},
		{`2001:10::`, true},
		{`2001:20::`, true},
		{`2001:db8::`, true},
		{`2002::`, true},
		{`2620:4f:8000::`, true},
		{`fc00::`, true},
		{`fe80::`, true},

		// https://datatracker.ietf.org/doc/html/rfc4291#section-2.2
		{`ABCD:EF01:2345:6789:ABCD:EF01:2345:6789`, true},
		{`2001:DB8:0:0:8:800:200C:417A`, true},
		{`2001:DB8:0:0:8:800:200C:417A`, true}, // a unicast address
		{`FF01:0:0:0:0:0:0:101`, true},         // a multicast address
		{`0:0:0:0:0:0:0:1`, true},              // the loopback address
		{`0:0:0:0:0:0:0:0`, true},              // the unspecified address
		{`2001:DB8::8:800:200C:417A`, true},    // a unicast address
		{`FF01::101`, true},                    // a multicast address
		{`::1`, true},                          // the loopback address
		//{`::`, true},                         // the unspecified address
		{`::`, nil},
		{`0:0:0:0:0:0:13.1.68.3`, true},
		{`0:0:0:0:0:FFFF:129.144.52.38`, true},
		{`::13.1.68.3`, true},
		{`::FFFF:129.144.52.38`, true},

		// https://datatracker.ietf.org/doc/html/rfc5952#section-1
		{`2001:db8:0:0:1:0:0:1`, true},
		{`2001:0db8:0:0:1:0:0:1`, true},
		{`2001:db8::1:0:0:1`, true},
		{`2001:db8::0:1:0:0:1`, true},
		{`2001:0db8::1:0:0:1`, true},
		{`2001:db8:0:0:1::1`, true},
		{`2001:db8:0000:0:1::1`, true},
		{`2001:DB8:0:0:1::1`, true},

		// https://datatracker.ietf.org/doc/html/rfc5952#section-2.1
		{`2001:db8:aaaa:bbbb:cccc:dddd:eeee:0001`, true},
		{`2001:db8:aaaa:bbbb:cccc:dddd:eeee:001`, true},
		{`2001:db8:aaaa:bbbb:cccc:dddd:eeee:01`, true},
		{`2001:db8:aaaa:bbbb:cccc:dddd:eeee:1`, true},

		// https://datatracker.ietf.org/doc/html/rfc5952#section-2.2
		{`2001:db8:aaaa:bbbb:cccc:dddd::1`, true},
		{`2001:db8:aaaa:bbbb:cccc:dddd:0:1`, true},
		{`2001:db8:0:0:0::1`, true},
		{`2001:db8:0:0::1`, true},
		{`2001:db8:0::1`, true},
		{`2001:db8::1`, true},
		{`2001:db8::aaaa:0:0:1`, true},
		{`2001:db8:0:0:aaaa::1`, true},

		// https://datatracker.ietf.org/doc/html/rfc5952#section-2.3
		{`2001:db8:aaaa:bbbb:cccc:dddd:eeee:aaaa`, true},
		{`2001:db8:aaaa:bbbb:cccc:dddd:eeee:AAAA`, true},
		{`2001:db8:aaaa:bbbb:cccc:dddd:eeee:AaAa`, true},

		// An IP address in URI host position must be bracketed unless it is IPv4.
		// https://www.rfc-editor.org/rfc/rfc3986#section-3.2.2
		// TODO: Implement this restriction, ideally without matching the `http://1080` prefix.
		//{`http://1080::8:800:200c:417a/path`, `1080::8:800:200c:417a`},

		{`foo.com:8080`, true},
		{`foo.com:8080/path`, true},
		{`test.foo.com`, true},
		{`test.foo.com/path`, true},
		{`test.foo.com/path/more/`, true},
		{`TEST.FOO.COM/PATH`, true},
		{`TEST.FÓO.COM/PÁTH`, true},
		{`foo.com/path_(more)`, true},
		{`foo.com/path_(even)_(more)`, true},
		{`foo.com/path_(more)/more`, true},
		{`foo.com/path_(more)/end)`, `foo.com/path_(more)/end`},
		{`www.foo.com`, true},
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
		{`foo@bar.com`, true},
		{`foo@sub.bar.com`, true},
		{`foo@bar.com bar@bar.com`, `foo@bar.com`},
		{`foo@bar.onion`, true},
		{`foo@中国.中国`, true},
		{`foo@test.bar.com`, true},
		{`FOO@TEST.BAR.COM`, true},
		{`foo@bar.com/path`, `foo@bar.com`},
		{`foo+test@bar.com`, true},
		{`foo+._%-@bar.com`, true},
	})
	doTest(t, "Strict2", Strict(), []testCase{
		{`http:// foo.com`, nil},
		{`foo.a`, nil},
		{`foo.com`, nil},
		{`foo.com/`, nil},
		{`1.1.1.1`, nil},
		{`3ffe:2a00:100:7031::1`, nil},
		{`test.foo.com:8080/path`, nil},
		{`foo@bar.com`, nil},

		// An IP address in URI host position must be bracketed unless it is IPv4.
		// https://www.rfc-editor.org/rfc/rfc3986#section-3.2.2
		// TODO: Implement this restriction, ideally without matching the `http://1080` prefix.
		//{`http://1080::8:800:200c:417a/path`, nil},
	})
}

func TestStrictMatchingSchemeError(t *testing.T) {
	for _, c := range []struct {
		exp     string
		wantErr bool
	}{
		{`http://`, false},
		{`https?://`, false},
		{`http://|mailto:`, false},
		{`http://(`, true},
	} {
		_, err := StrictMatchingScheme(c.exp)
		if c.wantErr && err == nil {
			t.Errorf(`StrictMatchingScheme("%s") did not error as expected`, c.exp)
		} else if !c.wantErr && err != nil {
			t.Errorf(`StrictMatchingScheme("%s") unexpectedly errored`, c.exp)
		}
	}
}

func TestStrictMatchingScheme(t *testing.T) {
	strictMatching, _ := StrictMatchingScheme("http://|ftps?://|mailto:")
	doTest(t, "StrictMatchingScheme", strictMatching, []testCase{
		{`foo.com`, nil},
		{`foo@bar.com`, nil},
		{`http://foo`, true},
		{`Http://foo`, true},
		{`https://foo`, nil},
		{`ftp://foo`, true},
		{`ftps://foo`, true},
		{`mailto:foo`, true},
		{`MAILTO:foo`, true},
		{`sms:123`, nil},
	})
}

func TestStrictMatchingSchemeAny(t *testing.T) {
	strictMatching, _ := StrictMatchingScheme(AnyScheme)
	doTest(t, "StrictMatchingScheme", strictMatching, []testCase{
		{`http://foo`, true},
		{`git+https://foo`, true},
		{`randomtexthttp://foo.bar/etc`, true},
		{`mailto:foo`, true},
	})
}

func bench(b *testing.B, re func() *regexp.Regexp, str string) {
	b.ReportAllocs()
	b.SetBytes(int64(len(str)))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			re().FindAllString(str, -1)
		}
	})
}

const inputNone = `
foo bar
yaml: "as well"
some more plaintext
which does not contain any urls.
`

const inputMany = `
foo bar http://foo.foo https://192.168.1.1/path
foo.com bitcoin:address ftp://
xmpp:foo@bar.com
`

func BenchmarkStrict_none(b *testing.B) {
	bench(b, Strict, inputNone)
}

func BenchmarkStrict_many(b *testing.B) {
	bench(b, Strict, inputMany)
}

func BenchmarkRelaxed_none(b *testing.B) {
	bench(b, Relaxed, inputNone)
}

func BenchmarkRelaxed_many(b *testing.B) {
	bench(b, Relaxed, inputMany)
}

var (
	rxMatchingScheme     *regexp.Regexp
	rxMatchingSchemeOnce sync.Once
)

func matchingScheme() *regexp.Regexp {
	rxMatchingSchemeOnce.Do(func() {
		rx, err := StrictMatchingScheme("https?://")
		if err != nil {
			panic(err)
		}
		rxMatchingScheme = rx
	})
	return rxMatchingScheme
}

func BenchmarkStrictMatchingScheme_none(b *testing.B) {
	bench(b, matchingScheme, inputNone)
}

func BenchmarkStrictMatchingScheme_many(b *testing.B) {
	bench(b, matchingScheme, inputMany)
}
