package main

import (
	"flag"
	"log"
	"os"

	"github.com/likexian/whois"
	"github.com/matti/gowhois"
)

func main() {
	flag.Parse()
	data := ""
	queryOrFile := flag.Arg(0)

	if _, err := os.Stat(queryOrFile); err == nil {
		bytes, err := os.ReadFile(queryOrFile)
		if err != nil {
			log.Panicln(err)
		}
		data = string(bytes)
	} else {
		var err error
		data, err = whois.Whois(queryOrFile)
		if err != nil {
			log.Panicln(err)
		}
	}

	//fmt.Println(data)

	gowhois.Parse(data)
}
