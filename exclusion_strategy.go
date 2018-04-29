package sieve

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type ExclusionStrategy interface {
	Name() string
	Check(v ...reflect.Value) bool
}

type excludeEqualField struct {
	fieldName string
}

func (excludeEqualField) Name() string {
	return "equal_filed"
}

func (e excludeEqualField) Check(v ...reflect.Value) bool {
	if len(v) < 2 {
		return false
	}
	if !v[0].CanInterface() {
		return false
	}
	kind := v[1].Type().Kind()
	if kind == reflect.Map {
		mapItemValue := v[1].MapIndex(reflect.ValueOf(e.fieldName))
		if mapItemValue.IsValid() && mapItemValue.CanInterface() {
			return reflect.DeepEqual(v[0].Interface(), mapItemValue.Interface())
		}
	} else if kind == reflect.Struct {
		fieldValue := v[1].FieldByName(e.fieldName)
		if fieldValue.IsValid() && fieldValue.CanInterface() {
			return reflect.DeepEqual(v[0].Interface(), fieldValue.Interface())
		}
	}
	return false
}

type excludeEqualValue struct {
	value string
}

func (excludeEqualValue) Name() string {
	return "equal_value"
}

func (e excludeEqualValue) Check(v ...reflect.Value) bool {
	if len(v) < 1 {
		return false
	}
	switch v[0].Type().Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(e.value, 10, 64)
		if err != nil {
			return false
		}
		return i == v[0].Int()
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(e.value, 10, 64)
		if err != nil {
			return false
		}
		return u == v[0].Uint()
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(e.value, 64)
		if err != nil {
			return false
		}
		return f == v[0].Float()
	case reflect.Bool:
		b, err := strconv.ParseBool(e.value)
		if err != nil {
			return false
		}
		return b == v[0].Bool()
	case reflect.String:
		return v[0].String() == e.value
	case reflect.Slice:
		if v[0].Elem().Kind() == reflect.Uint8 {
			return bytes.Equal(v[0].Bytes(), []byte(e.value))
		}
	case reflect.Struct:
		{
			if !v[0].CanInterface() {
				return false
			}
			switch rv := v[0].Interface().(type) {
			case time.Time:
				rv.String()
				if t, err := time.Parse(time.RFC3339, e.value); err == nil {
					return rv.Equal(t)
				}
				if timeStamp, err := strconv.ParseInt(e.value, 10, 64); err == nil {
					return rv.Equal(time.Unix(timeStamp, 0))
				}
			case fmt.Stringer:
				return rv.String() == e.value
			}
		}
	}
	return false
}
