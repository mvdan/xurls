/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

var regexLink = regexp.MustCompile(
	`(([^\s'"<>\(\)]+:(//)?|(http|ftp|www)[^.]*\.)[^\s'"<>\(\)]*|[^\s'"<>\(\)]+\.(com|org|net|edu|info)(/[^\s'"<>\(\)]*)?)[^.,;:\s'"<>\(\)]`)

func FindString(s string) string {
	return regexLink.FindString(s)
}

func FindAllString(s string) []string {
	return regexLink.FindAllString(s, -1)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		urls := FindAllString(line)
		if urls == nil {
			continue
		}
		for _, url := range urls {
			fmt.Println(url)
		}
	}
}
