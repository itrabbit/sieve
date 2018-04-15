package sieve

import (
	"sort"
)

type sieve struct {
	v      interface{}
	scopes []string
}

func (s sieve) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.v, &s)
}

func (s sieve) HasScope(scope string) bool {
	if scope == "*" {
		return true
	}
	if len(s.scopes) > 0 {
		i := sort.SearchStrings(s.scopes, scope)
		return i >= 0 && i < len(s.scopes) && s.scopes[i] == scope
	}
	return false
}

func (s sieve) HasAnyScope(scopes ...string) bool {
	for _, scope := range scopes {
		if s.HasScope(scope) {
			return true
		}
	}
	return false
}

func Sieve(v interface{}, scopes ...string) interface{} {
	if len(scopes) > 1 {
		sort.Strings(scopes)
	}
	return &sieve{v, scopes}
}
