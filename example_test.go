// Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package xurls_test

import (
	"fmt"

	"mvdan.cc/xurls"
)

func Example() {
	urlsRe := xurls.Relaxed()
	fmt.Println(urlsRe.FindString("Do gophers live in http://golang.org?"))
	fmt.Println(urlsRe.FindAllString("foo.com is http://foo.com/.", -1))
	// Output:
	// http://golang.org
	// [foo.com http://foo.com/]
}
