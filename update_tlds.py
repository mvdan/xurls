#!/usr/bin/env python3

import urllib.request

TLDSURL = "https://data.iana.org/TLD/tlds-alpha-by-domain.txt"
OUTFILE = "tlds.go"


def urlget(url):
    response = urllib.request.urlopen(url)
    data = response.read()
    text = data.decode('utf-8')
    return text

def main():
    text = urlget(TLDSURL)
    lines = text.splitlines()
    tlds = [line.lower() for line in lines if line[0] != '#']
    # Reversed so that the longest match first
    tlds.sort(reverse=True)
    with open(OUTFILE, mode='w+') as f:
        f.write("""\
/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package xurls

var tlds = []string{
""")
        for tld in tlds:
            f.write("\t\"%s\",\n" % tld)
        f.write("""\
}
""")

if __name__ == '__main__':
    main()
