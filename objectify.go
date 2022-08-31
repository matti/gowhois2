package gowhois

import (
	"log"
	"regexp"
	"strings"
)

func objectify(lines []string) []Entity {
	var objects []Entity

	object := Entity{}
	r := regexp.MustCompile(`^([^:]+):\s+\S`)

	for _, line := range lines {
		if line == "" {
			if object.Kind != "" {
				objects = append(objects, object)
				object = Entity{}
			}
		} else if strings.HasPrefix(line, "#") {
		} else if strings.HasPrefix(line, ";;") {
		} else if strings.HasPrefix(line, "%") {
		} else if r.MatchString(line) {
			parts := strings.SplitN(line, ":", 2)
			key := parts[0]
			value := parts[1]

			if object.Kind == "" {
				object.Kind = key
			}

			var lastPair Pair
			if len(object.Pairs) > 0 {
				lastPair = object.Pairs[len(object.Pairs)-1]
			}
			if lastPair.Key == key {
				appended := lastPair.Value + "\n" + strings.TrimSpace(value)
				object.Pairs[len(object.Pairs)-1] = Pair{
					Key:   key,
					Value: appended,
				}
				log.Println(key, appended)
			} else {
				object.Pairs = append(object.Pairs, Pair{
					Key:   key,
					Value: strings.TrimSpace(value),
				})
			}
		} else if len(object.Pairs) > 0 && strings.HasPrefix(line, " ") {
			lastPair := object.Pairs[len(object.Pairs)-1]
			lastPair.Value = lastPair.Value + "\n" + strings.TrimSpace(line)
			object.Pairs[len(object.Pairs)-1] = lastPair
		} else if strings.HasSuffix(strings.TrimSpace(line), ":") {
			if len(object.Pairs) > 0 {
				lastPair := object.Pairs[len(object.Pairs)-1]
				parts := strings.SplitN(line, ":", 2)
				key := parts[0]

				if lastPair.Key == key {
					appended := lastPair.Value + "\n"
					object.Pairs[len(object.Pairs)-1] = Pair{
						Key:   key,
						Value: appended,
					}
				}
			} else {
				// "OriginAS:   "
			}
		} else {
			log.Fatalln("unknown line", "'"+line+"'")
		}
	}

	return objects
}
