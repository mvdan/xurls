/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mvdan/xurls"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		urls := xurls.All.FindAllString(line, -1)
		if urls == nil {
			continue
		}
		for _, url := range urls {
			fmt.Println(url)
		}
	}
}
