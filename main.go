package main

import (
	"flag"
	"fmt"
)

var ver = flag.String("v", "", "the go version you want to download")
var url = flag.String("u", "", "the url of the go version you want to download")
var helpMsg = "Usage: go-update-install -v <version-number>\n" +
	"   or: go-update-install -u <url of binary to install"

func try(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	flag.Parse()
	if *url == "" {
		if *ver != "" {
			*url = "go" + *ver + ".linux-amd64.tar.gz"
		} else {
			fmt.Println(helpMsg)
			return
		}
	}
	try(download(*url))
	uninstallOld
	installNew
}
