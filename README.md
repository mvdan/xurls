# xurls

[![GoDoc](https://godoc.org/github.com/mvdan/xurls?status.svg)](https://godoc.org/github.com/mvdan/xurls)
[![Travis](https://travis-ci.org/mvdan/xurls.svg?branch=master)](https://travis-ci.org/mvdan/xurls)

Extract urls from text using regular expressions.

	go get github.com/mvdan/xurls

#### cmd/xurls

	go get github.com/mvdan/xurls/cmd/xurls

```shell
$ echo "Do gophers live in http://golang.org?" | xurls
http://golang.org
```
