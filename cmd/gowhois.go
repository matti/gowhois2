package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/likexian/whois"
	"github.com/matti/gowhois2"
)

func main() {
	flag.Parse()
	queryOrFile := flag.Arg(0)

	if stdinFileInfo, err := os.Stdin.Stat(); err == nil {
		if (stdinFileInfo.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				input := scanner.Text()
				data, err := whois.Whois(input)
				if err != nil {
					panic(err)
				}
				log.Println(data)
				fmt.Println(gowhois2.Parse(data))
			}
		}
	}
	var data string
	if _, err := os.Stat(queryOrFile); err == nil {
		bytes, err := os.ReadFile(queryOrFile)
		if err != nil {
			panic(err)
		}
		data = string(bytes)
	} else {
		var err error
		data, err = whois.Whois(queryOrFile)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(gowhois2.Parse(data))
}
