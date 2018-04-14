package sieve

import (
	"sort"
)

// Сито
type sieve struct {
	v      interface{}
	groups []string
}

// Сериаизация в JSON
func (s sieve) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.v, &s)
}

// Содержит ли группу
func (s sieve) HasGroup(group string) bool {
	if group == "*" {
		return true
	}
	if len(s.groups) > 0 {
		i := sort.SearchStrings(s.groups, group)
		return i >= 0 && i < len(s.groups) && s.groups[i] == group
	}
	return false
}

// Содержит любую из групп
func (s sieve) HasAnyGroup(groups ...string) bool {
	for _, group := range groups {
		if s.HasGroup(group) {
			return true
		}
	}
	return false
}

// Создаем набор из данных для сито
func Sieve(v interface{}, groups ...string) interface{} {
	if len(groups) > 1 {
		sort.Strings(groups)
	}
	return &sieve{v, groups}
}
