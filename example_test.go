// Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package xurls_test

import (
	"fmt"

	"mvdan.cc/xurls/v2"
)

func Example() {
	rx := xurls.Relaxed()
	fmt.Println(rx.FindString("Do gophers live in http://golang.org?"))
	fmt.Println(rx.FindAllString("foo.com is http://foo.com/.", -1))
	// Output:
	// http://golang.org
	// [foo.com http://foo.com/]
}
