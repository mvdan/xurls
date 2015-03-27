# xurls

Extract urls from plain text using regular expressions.

	go get github.com/mvdan/xurls

Example usage:

```go
package main

import "github.com/mvdan/xurls"

func main() {
        xurls.All.FindString("Do gophers live in golang.org?")
        // "golang.org"
        xurls.All.FindAllString("foo.com is http://foo.com/.", -1)
        // ["foo.com", "http://foo.com/"]
        xurls.AllStrict.FindAllString("foo.com is http://foo.com/.", -1)
        // ["http://foo.com/"]
}
```

### Command-line utilities

#### xurls

Reads text and prints one url per line.

	go get github.com/mvdan/xurls/cmd/xurls

```shell
$ echo "Do gophers live in golang.org?" | xurls
golang.org
```

* **-s** only match urls with scheme (strict)
