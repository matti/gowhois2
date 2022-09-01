package gowhois

import (
	"encoding/json"
	"regexp"
	"strings"
)

type Entity struct {
	Name       string              `json:"name"`
	Related    string              `json:"related"`
	Properties map[string][]string `json:"properties"`
}
type Registry struct {
	Name     string   `json:"name"`
	Entities []Entity `json:"entities"`
}

type Response struct {
	Registries []Registry `json:"registries"`
}

func (r Response) String() string {
	sJSON, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(sJSON)
}

var registryLines = map[string]string{
	"% IANA WHOIS server": "iana",
	"# whois.arin.net":    "arin",
	"# whois.ripe.net":    "ripe",
	"% [whois.apnic.net]": "apnic",
	"# ARIN WHOIS data and services are subject to the Terms of Use": "arin",
	"% This is the RIPE Database query service.":                     "ripe",
}

func Parse(data string) Response {
	lines := strings.Split(data, "\n")

	var response Response

	var registry Registry
	var entity Entity
	var key string
	var related string

	keyValueLineRegExp := regexp.MustCompile(`^([^:]+):`)
	relatedLineRegExp := regexp.MustCompile(`^% Information related to '([^']+)'`)
	for _, line := range lines {
		// Registry needs to be checked first before ignoring comments because they start with comment char
		if r := registryLines[line]; r != "" {
			// if new registry is found, append to registries first
			if registry.Name != "" && len(registry.Entities) > 0 {
				response.Registries = append(response.Registries, registry)
			}
			registry = Registry{
				Name:     r,
				Entities: []Entity{},
			}
			continue
		}

		if line == "" {
			// Empty line stops entity
			if registry.Name != "" && len(entity.Properties) > 0 {
				registry.Entities = append(registry.Entities, entity)
				entity = Entity{}
			}
			continue
		}

		if matches := relatedLineRegExp.FindStringSubmatch(line); len(matches) > 1 {
			related = matches[1]
		}

		if strings.HasPrefix(line, "#") ||
			strings.HasPrefix(line, ";;") ||
			strings.HasPrefix(line, "%") {
			continue
		}

		switch registry.Name {
		case "":
			panic("failed to detect registry")
		case "iana", "ripe", "arin", "apnic":
		default:
			panic("unknown registry: '" + registry.Name + "'")
		}

		if keyValueLineRegExp.MatchString(line) {
			parts := strings.SplitN(line, ":", 2)
			key = parts[0]
			// Value can be empty when line "OriginAS:"
			value := parts[1]

			// if entity has no name, use the first property key as the name
			if entity.Name == "" {
				entity.Name = key
				entity.Related = related
				entity.Properties = make(map[string][]string)
			}

			entity.Properties[key] = append(
				entity.Properties[key],
				strings.TrimSpace(value),
			)

			continue
		}

		if len(entity.Properties[key]) > 0 && strings.HasPrefix(line, " ") {
			// concat continuation, eg "             Somestreet 6" to the previous property
			value := entity.Properties[key][len(entity.Properties[key])-1] + "\n" + strings.TrimSpace(line)
			entity.Properties[key][len(entity.Properties[key])-1] = strings.TrimSpace(value)

			continue
		}

		panic("unknown line '" + line + "'")
	}

	// arin.txt has arin identifier again at the end, do not add a new registry without entities
	if len(registry.Entities) > 0 {
		response.Registries = append(response.Registries, registry)
	}

	return response
}
