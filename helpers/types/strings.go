package types

import (
	"reflect"
	"strconv"
)

func StringToType(t reflect.Type, v string) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.Int:
		i, err := strconv.Atoi(v)
		return reflect.ValueOf(i), err
	case reflect.Int8:
		i, err := strconv.ParseInt(v, 10, 8)
		return reflect.ValueOf(int8(i)), err
	case reflect.Int16:
		i, err := strconv.ParseInt(v, 10, 16)
		return reflect.ValueOf(int16(i)), err
	case reflect.Int32:
		i, err := strconv.ParseInt(v, 10, 32)
		return reflect.ValueOf(int32(i)), err
	case reflect.Int64:
		i, err := strconv.ParseInt(v, 10, 64)
		return reflect.ValueOf(int64(i)), err

	case reflect.Uint:
		i, err := strconv.Atoi(v)
		return reflect.ValueOf(uint(i)), err
	case reflect.Uint8:
		i, err := strconv.ParseUint(v, 10, 8)
		return reflect.ValueOf(uint8(i)), err
	case reflect.Uint16:
		i, err := strconv.ParseUint(v, 10, 16)
		return reflect.ValueOf(uint16(i)), err
	case reflect.Uint32:
		i, err := strconv.ParseUint(v, 10, 32)
		return reflect.ValueOf(uint32(i)), err
	case reflect.Uint64:
		i, err := strconv.ParseUint(v, 10, 64)
		return reflect.ValueOf(uint64(i)), err
	}

	return reflect.ValueOf(v), nil
}
