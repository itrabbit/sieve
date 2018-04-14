package sieve

import (
	"sort"
	"strings"
)

type tagOptions struct {
	scopes            []string // s
	exportKeys        []string // ek
	excludeEqualField string   // eef
}

func parseTag(tag string) (opts tagOptions) {
	for _, param := range strings.Split(tag, ";") {
		param = strings.TrimSpace(param)
		if len(param) < 3 {
			continue
		}
		if strings.HasPrefix(param, "s:") {
			opts.scopes = strings.Split(param[2:], ",")
			continue
		}
		if strings.HasPrefix(param, "ek:") {
			opts.exportKeys = strings.Split(param[3:], ",")
			sort.Strings(opts.exportKeys)
			continue
		}
		if strings.HasPrefix(param, "eef:") {
			opts.excludeEqualField = strings.TrimSpace(param[4:])
			continue
		}
	}
	return
}
