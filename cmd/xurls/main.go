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
	re := xurls.All
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		matches := re.FindAllString(word, -1)
		for _, match := range matches {
			fmt.Println(match)
		}
	}
}
