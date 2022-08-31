package gowhois

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Pair struct {
	Key   string
	Value string
}
type Entity struct {
	Kind  string
	Pairs []Pair
}
type Record struct {
	Registry string
	Entities []Entity
}

func Parse(data string) {
	lines := strings.Split(data, "\n")

	registry := ""
	for _, line := range lines {
		switch line {
		case "% [whois.apnic.net]":
			registry = "apnic"
		case "# ARIN WHOIS data and services are subject to the Terms of Use":
			registry = "arin"
		case "% This is the RIPE Database query service.":
			registry = "ripe"
		}
	}

	var entities []Entity
	switch registry {
	case "ripe", "arin", "apnic":
		entities = objectify(lines)
	default:
		log.Fatalln("unknown registry:", "'"+registry+"'")
	}

	record := &Record{
		Registry: registry,
		Entities: entities,
	}
	sJSON, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(string(sJSON))
}
