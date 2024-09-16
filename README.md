# xurls

> A modified version of [xurls](github.com/mvdan/xurls), with a focus on matching content within user messages that would be clickable in web browsers.
>
> In short, this stops a bug with [Fossabot on YouTube](https://fossabot.com), where it accidentally matches multiple YouTube emotes chained together due to `::` being in the string.
>
> This seems like desired/intentional behaviour from `xurls`, hence the fork and not contributing back to the package. The package is correct, in relaxed mode, these ARE valid ipv6 addresses! However, this causes problems for this specific use case.
>
> Please contribute back and star the original xurls package!

[![Go Reference](https://pkg.go.dev/badge/github.com/aidenwallis/xurls.svg)](https://pkg.go.dev/github.com/aidenwallis/xurls)

Extract urls from text using regular expressions. Requires Go 1.22 or later.

```go
import "github.com/aidenwallis/xurls"

func main() {
	rxRelaxed := xurls.Relaxed()
	rxRelaxed.FindString("Do gophers live in golang.org?")  // "golang.org"
	rxRelaxed.FindString("This string does not have a URL") // ""

	rxStrict := xurls.Strict()
	rxStrict.FindAllString("must have scheme: http://foo.com/.", -1) // []string{"http://foo.com/"}
	rxStrict.FindAllString("no scheme, no match: foo.com", -1)       // []string{}
}
```

Since API is centered around [regexp.Regexp](https://golang.org/pkg/regexp/#Regexp),
many other methods are available, such as finding the [byte indexes](https://golang.org/pkg/regexp/#Regexp.FindAllIndex)
for all matches.

The regular expressions are compiled when the API is first called.
Any subsequent calls will use the same regular expression pointers.

#### cmd/xurls

To install the tool globally:

	go install github.com/aidenwallis/xurls/cmd/xurls@latest

```shell
$ echo "Do gophers live in http://golang.org?" | xurls
http://golang.org
```
