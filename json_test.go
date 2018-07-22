package sieve

import (
	"bytes"
	"encoding/json"
	"testing"
)

type ObjectOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Object struct {
	idx       uint64
	Name      string         `json:"name" sieve:"s:*"`
	FullName  string         `json:"full_name" sieve:"s:private"`
	CreatedAt uint64         `json:"created_at"`
	UpdatedAt uint64         `json:"updated_at" sieve:"eef:CreatedAt"`
	Options   []ObjectOption `json:"options" sieve:"k:Name"`
}

type Nesting struct {
	ObjectOption
	CreatedAt uint64 `json:"created_at"`
	UpdatedAt uint64 `json:"updated_at" sieve:"eef:CreatedAt"`
}

type ObjectN struct {
	a         ObjectOption
	CreatedAt uint64 `json:"created_at"`
	UpdatedAt uint64 `json:"updated_at" sieve:"eef:CreatedAt"`
}

func TestSieveMarshalJSON_UnexportedFields(t *testing.T) {
	obj := ObjectN{
		a: ObjectOption{
			Name:  "One",
			Value: "Value",
		},
		CreatedAt: 100,
		UpdatedAt: 100,
	}
	b, err := json.Marshal(Sieve(&obj, "public"))
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Contains(b, []byte("name")) {
		t.Fail()
		return
	}
	if bytes.Contains(b, []byte("value")) {
		t.Fail()
		return
	}
}

func TestSieveMarshalJSON_Nested(t *testing.T) {
	obj := Nesting{
		ObjectOption: ObjectOption{
			Name:  "One",
			Value: "Value",
		},
		CreatedAt: 100,
		UpdatedAt: 100,
	}
	b, err := json.Marshal(Sieve(&obj, "public"))
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Contains(b, []byte("name")) {
		t.Fail()
		return
	}
	if !bytes.Contains(b, []byte("value")) {
		t.Fail()
		return
	}
}

func TestSieveMarshalJSON_Scopes(t *testing.T) {
	obj := Object{
		idx:      100,
		Name:     "One",
		FullName: "Full",
	}
	b, err := json.Marshal(Sieve(&obj, "public"))
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Contains(b, []byte("full_name")) {
		t.Fail()
		return
	}
	b, err = json.Marshal(Sieve(&obj, "private"))
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Contains(b, []byte("full_name")) {
		t.Fail()
		return
	}
}

func TestSieveMarshalJSON_ExcludeEqualField(t *testing.T) {
	obj := Object{
		CreatedAt: 100,
		UpdatedAt: 100,
	}
	b, err := json.Marshal(Sieve(&obj))
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Contains(b, []byte("updated_at")) {
		t.Fail()
		return
	}
	obj.UpdatedAt = 120
	b, err = json.Marshal(Sieve(&obj))
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Contains(b, []byte("updated_at")) {
		t.Fail()
		return
	}
}

func TestSieveMarshalJSON_ExportKeys(t *testing.T) {
	obj := Object{
		Options: []ObjectOption{
			ObjectOption{"1", "First"},
			ObjectOption{"2", "Two"},
			ObjectOption{"3", "Three"},
		},
	}
	b, err := json.Marshal(Sieve(&obj))
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Contains(b, []byte("[\"1\",\"2\",\"3\"]")) {
		t.Fail()
		return
	}
}
