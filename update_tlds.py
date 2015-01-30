#!/usr/bin/env python3

import re
import urllib.request

def urlget(url):
    response = urllib.request.urlopen(url)
    data = response.read()
    text = data.decode('utf-8')
    return text

def main():
    tlds = set()

    regex = re.compile(r"^([^#]+)$")
    data = urlget("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
    for line in data.splitlines():
        m = regex.search(line)
        if not m:
            continue
        tld = m.group(1)
        tlds.add(tld.lower())

    regex = re.compile(r"(^([^/.]+)$|^// (xn--[^\s]+)[\s$])")
    data = urlget("https://publicsuffix.org/list/effective_tld_names.dat")
    for line in data.splitlines():
        m = regex.search(line)
        if not m:
            continue
        tld = m.group(2) or m.group(3)
        tlds.add(tld)

    # Reversed so that the longest match first
    tldslist = [t for t in tlds]
    tldslist.sort()
    with open("tlds.go", mode='w+') as f:
        f.write("""\
/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package xurls

var TLDs = []string{
""")
        for tld in tldslist:
            f.write("\t`%s`,\n" % tld)
        f.write("""\
}
""")

if __name__ == '__main__':
    main()
