package sieve

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Option struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Object struct {
	idx       uint64
	Name      string   `json:"name" sieve:"g:*"`
	FullName  string   `json:"full_name" sieve:"g:private"`
	CreatedAt uint64   `json:"created_at"`
	UpdatedAt uint64   `json:"updated_at" sieve:"eef:CreatedAt"`
	Options   []Option `json:"options" sieve:"ek:Name"`
}

func TestSievePrivateMarshalJSON(t *testing.T) {
	obj := Object{
		100,
		"One",
		"Full",
		100,
		100,
		[]Option{
			Option{"1", "Cool"},
			Option{"2", "Great"},
			Option{"3", "Best"},
		},
	}
	b, err := json.Marshal(Sieve(&obj, "public"))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}
