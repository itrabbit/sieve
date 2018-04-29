package sieve

import (
	"sort"
	"strings"
)

func parseTag(tag string) Options {
	opts := &options{}
	for _, param := range strings.Split(tag, ";") {
		param = strings.TrimSpace(param)
		if len(param) < 3 {
			continue
		}
		chunks := strings.SplitN(param, ":", 2)
		if len(chunks) < 2 {
			if chunks[0] == "e.any" {
				opts.SetIsAnyExclusion(true)
			}
			continue
		}
		switch strings.TrimSpace(chunks[0]) {
		case "s", "scopes":
			opts.scopes = strings.Split(chunks[1], ",")
		case "k", "exportKeys":
			opts.exportKeys = strings.Split(chunks[1], ",")
			sort.Strings(opts.exportKeys)
		case "ef", "eef", "excludeEqualField":
			opts.PutExclusionStrategy(&excludeEqualField{strings.TrimSpace(chunks[1])})
		case "ev", "eev", "excludeEqualValue":
			opts.PutExclusionStrategy(&excludeEqualValue{chunks[1]})
		}
	}
	return opts
}
