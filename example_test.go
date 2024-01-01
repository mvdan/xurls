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

func ExampleStrictMatchingScheme() {
	rx, err := xurls.StrictMatchingScheme(`https?://`)
	if err != nil {
		panic(err)
	}
	fmt.Println(rx.FindAllString("Download binaries via https://foo.com/dl or ftps://foo.com/dl", -1))
	// Output:
	// [https://foo.com/dl]
}

func Example_filterEmails() {
	s := "Email dev@foo.com about any issues with foo.com or https://foo.com/dl"
	rx := xurls.Relaxed()
	idxEmail := rx.SubexpIndex("relaxedEmail")
	for _, match := range rx.FindAllStringSubmatch(s, -1) {
		if match[idxEmail] != "" {
			continue // skip lone email addresses
		}
		fmt.Println(match[0])
	}
	// Output:
	// foo.com
	// https://foo.com/dl
}
