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
	MarshalSieveJSON(opts Options) ([]byte, error)
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
		return i.MarshalSieveJSON(BuildOptions(s.scopes, nil))
	}
	if _, ok := v.(json.Marshaler); ok {
		return json.Marshal(v)
	}
	if k == reflect.Struct {
		obj, err := convertValueToMap(reflect.Indirect(reflect.ValueOf(v)), s, nil)
		if err != nil {
			return nil, err
		}
		return json.Marshal(obj)
	}
	if k == reflect.Map {
		obj, err := bustValueMap(reflect.Indirect(reflect.ValueOf(v)), s, nil)
		if err != nil {
			return nil, err
		}
		return json.Marshal(obj)
	}
	return json.Marshal(v)
}

func bustValue(val reflect.Value, s *sieve, exportKeys []string) (interface{}, error) {
	if !val.CanInterface() {
		return nil, nil
	}
	if s == nil {
		return val.Interface(), nil
	}
	if i, ok := val.Interface().(json.Marshaler); ok {
		return i, nil
	}
	if i, ok := val.Interface().(MarshalerJSON); ok {
		b, err := i.MarshalSieveJSON(BuildOptions(s.scopes, exportKeys))
		if err != nil {
			return nil, err
		}
		return json.RawMessage(b), nil
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
	return val.Interface(), nil
}

func bustValueSlice(val reflect.Value, s *sieve, exportKeys []string) (interface{}, error) {
	if !val.CanInterface() {
		return nil, nil
	}
	if s == nil || val.Len() == 0 {
		return val.Interface(), nil
	}
	if exportKeys != nil && len(exportKeys) > 0 {
		list := make([]interface{}, 0)
		for index := 0; index < val.Len(); index++ {
			i, err := bustValue(reflect.Indirect(val.Index(index)), s, exportKeys)
			if err != nil {
				return nil, err
			}
			if i != nil {
				list = append(list, i)
			}
		}
		return list, nil
	}
	list := make([]interface{}, val.Len(), val.Len())
	for index := 0; index < val.Len(); index++ {
		item, err := bustValue(reflect.Indirect(val.Index(index)), s, exportKeys)
		if err != nil {
			return nil, err
		}
		list[index] = item
	}
	return list, nil
}

func bustValueMap(val reflect.Value, s *sieve, exportKeys []string) (interface{}, error) {
	if !val.CanInterface() {
		return nil, nil
	}
	if s == nil || val.Len() == 0 {
		return val.Interface(), nil
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
		obj, err := bustValue(val.MapIndex(key), s, nil)
		if err != nil {
			return nil, err
		}
		m[keyStr] = obj
	}
	return m, nil
}

func convertValueToMap(val reflect.Value, s *sieve, exportKeys []string) (interface{}, error) {
	if !val.CanInterface() {
		return nil, nil
	}
	if s == nil || !val.IsValid() {
		return val.Interface(), nil
	}
	t, exporting := val.Type(), false
	if count := len(exportKeys); count > 0 {
		exporting = true
		if count == 1 {
			if s, ok := t.FieldByName(exportKeys[0]); ok && !s.Anonymous {
				if c := val.FieldByName(exportKeys[0]); c.CanInterface() {
					return reflect.Indirect(c).Interface(), nil
				}
			}
			return nil, nil
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
		if opts.HasScopes() {
			if !s.HasAnyScope(opts.Scopes()...) {
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
		if opts.HasExclusions() {
			if opts.CheckByExclusions(fieldValue, val) {
				continue
			}
		}
		obj, err := bustValue(fieldValue, s, opts.ExportKeys())
		if err != nil {
			return nil, err
		}
		m[fieldName] = obj
	}
	return m, nil
}
