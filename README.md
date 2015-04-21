# xurls

Extract urls from plain text using regular expressions.

	go get github.com/mvdan/xurls

Example usage:

```go
import "github.com/mvdan/xurls"

func main() {
        xurls.Relaxed.FindString("Do gophers live in golang.org?")
        // "golang.org"
        xurls.Relaxed.FindAllString("foo.com is http://foo.com/.", -1)
        // ["foo.com", "http://foo.com/"]
        xurls.Strict.FindAllString("foo.com is http://foo.com/.", -1)
        // ["http://foo.com/"]
}
```

This is **not** a URL validation library. Extracted urls may well not be
valid.

### Command-line utilities

#### xurls

Reads text and prints one url per line.

	go get github.com/mvdan/xurls/cmd/xurls

```shell
$ echo "Do gophers live in http://golang.org?" | xurls
http://golang.org
```
