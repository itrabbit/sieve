package sieve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type H map[string]interface{}

type MarshalerJSON interface {
	MarshalSieveJSON(scopes []string, exportKeys []string) ([]byte, error)
}

func marshalJSON(v interface{}, s *sieve) ([]byte, error) {
	if s == nil {
		return json.Marshal(v)
	}
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	k := t.Kind()
	if k == reflect.Slice || k == reflect.Array {
		empty, buff := true, bytes.Buffer{}
		buff.WriteByte('[')
		if val := reflect.Indirect(reflect.ValueOf(v)); !val.IsNil() && val.IsValid() {
			for i := 0; i < val.Len(); i++ {
				b, err := marshalJSON(val.Index(i).Interface(), s)
				if err != nil {
					return nil, err
				}
				if b != nil && len(b) > 0 {
					if !empty {
						buff.WriteByte(',')
					} else {
						empty = true
					}
					buff.Write(b)
				}
			}
		}
		buff.WriteByte(']')
		return buff.Bytes(), nil
	}
	if i, ok := v.(MarshalerJSON); ok {
		if i == nil {
			return nil, nil
		}
		return i.MarshalSieveJSON(s.scopes, nil)
	}
	if _, ok := v.(json.Marshaler); ok {
		return json.Marshal(v)
	}
	if k == reflect.Struct {
		return json.Marshal(convertValueToMap(reflect.Indirect(reflect.ValueOf(v)), s, nil))
	}
	if k == reflect.Map {
		return json.Marshal(bustValueMap(reflect.Indirect(reflect.ValueOf(v)), s, nil))
	}
	return json.Marshal(v)
}

func bustValue(val reflect.Value, s *sieve, exportKeys []string) interface{} {
	if !val.CanInterface() {
		return nil
	}
	if s == nil {
		return val.Interface()
	}
	if i, ok := val.Interface().(json.Marshaler); ok {
		return i
	}
	if i, ok := val.Interface().(MarshalerJSON); ok {
		b, err := i.MarshalSieveJSON(s.scopes, exportKeys)
		if err != nil {
			return nil
		}
		return json.RawMessage(b)
	}
	kind := val.Kind()
	if kind == reflect.Array || kind == reflect.Slice {
		return bustValueSlice(val, s, exportKeys)
	}
	if kind == reflect.Map {
		return bustValueMap(val, s, exportKeys)
	}
	if kind == reflect.Struct {
		return convertValueToMap(val, s, exportKeys)
	}
	return val.Interface()
}

func bustValueSlice(val reflect.Value, s *sieve, exportKeys []string) interface{} {
	if !val.CanInterface() {
		return nil
	}
	if s == nil || val.Len() == 0 {
		return val.Interface()
	}
	if exportKeys != nil && len(exportKeys) > 0 {
		list := make([]interface{}, 0)
		for index := 0; index < val.Len(); index++ {
			i := bustValue(reflect.Indirect(val.Index(index)), s, exportKeys)
			if i != nil {
				list = append(list, i)
			}
		}
		return list
	}
	list := make([]interface{}, val.Len(), val.Len())
	for index := 0; index < val.Len(); index++ {
		list[index] = bustValue(reflect.Indirect(val.Index(index)), s, exportKeys)
	}
	return list
}

func bustValueMap(val reflect.Value, s *sieve, exportKeys []string) interface{} {
	if !val.CanInterface() {
		return nil
	}
	if s == nil || val.Len() == 0 {
		return val.Interface()
	}
	m, exporting, oneKey := make(H), len(exportKeys) > 0, len(exportKeys) == 1
	for _, key := range val.MapKeys() {
		if !key.CanInterface() {
			continue
		}
		keyStr := strings.TrimSpace(fmt.Sprint(key.Interface()))
		if len(keyStr) < 1 {
			continue
		}
		if exporting {
			idx := sort.SearchStrings(exportKeys, keyStr)
			if idx < 0 || idx >= len(exportKeys) {
				continue
			}
			if exportKeys[idx] != keyStr {
				continue
			}
			if oneKey {
				return bustValue(val.MapIndex(key), s, nil)
			}
		}
		m[keyStr] = bustValue(val.MapIndex(key), s, nil)
	}
	return m
}

func convertValueToMap(val reflect.Value, s *sieve, exportKeys []string) interface{} {
	if !val.CanInterface() {
		return nil
	}
	if s == nil || !val.IsValid() {
		return val.Interface()
	}
	t, exporting := val.Type(), false
	if count := len(exportKeys); count > 0 {
		exporting = true
		if count == 1 {
			if s, ok := t.FieldByName(exportKeys[0]); ok && !s.Anonymous {
				if c := val.FieldByName(exportKeys[0]); c.CanInterface() {
					return reflect.Indirect(c).Interface()
				}
			}
			return nil
		}
	}
	m := make(H)
	for index := 0; index < t.NumField(); index++ {
		field := t.Field(index)
		if field.Anonymous {
			continue
		}
		if exporting {
			idx := sort.SearchStrings(exportKeys, field.Name)
			if idx < 0 || idx >= len(exportKeys) {
				continue
			}
			if exportKeys[idx] != field.Name {
				continue
			}
		}
		fieldName, omitempty := field.Name, false
		if tag, ok := field.Tag.Lookup("json"); ok {
			if idx := strings.Index(tag, ","); idx != -1 {
				if name := strings.TrimSpace(tag[:idx]); len(name) > 0 {
					fieldName = name
				}
				omitempty = strings.Contains(tag[idx+1:], "omitempty")
			}
			if tag = strings.TrimSpace(tag); len(tag) > 0 {
				fieldName = tag
			}
		}
		opts := parseTag(field.Tag.Get("sieve"))
		if opts.scopes != nil && len(opts.scopes) > 0 {
			if !s.HasAnyScope(opts.scopes...) {
				continue
			}
		}
		fieldValue := reflect.Indirect(val.Field(index))
		if !fieldValue.CanInterface() {
			continue
		}
		if omitempty && isEmptyValue(fieldValue) {
			continue
		}
		if len(opts.excludeEqualField) > 0 {
			c := reflect.Indirect(val.FieldByName(opts.excludeEqualField))
			if c.CanInterface() && c.Type().Kind() == fieldValue.Type().Kind() {
				if reflect.DeepEqual(c.Interface(), fieldValue.Interface()) {
					continue
				}
			}
		}
		m[fieldName] = bustValue(fieldValue, s, opts.exportKeys)
	}
	return m
}
