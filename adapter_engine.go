package adapter

import (
	"fmt"
	"reflect"

	"github.com/vingarcia/tagmapper/tags"
)

// TagDecoder is the adapter that allows the Decode function to get values
// from any data source.
//
// The TagDecoder will then receive a tagValue and a reflect.Kind and return
// a value that is compatible with that type, e.g.:
//
// Suppose a decoder that reads from os.Getenv() and it receives:
// tagValue = "NUM_RETRIES" and fieldKind = reflect.Int, then it should
// return strconv.Atoi(os.Getenv("NUM_RETRIES"))
//
// The decoder will then proceed to check if the type matches the expected
// and then assign it to the struct field that contained the previous specified
// tagValue.
type TagDecoder interface {
	DecodeField(info Field) (interface{}, error)
}

type Field struct {
	Tags map[string]string
	Name string
	Kind reflect.Kind
	Type reflect.Type
}

func Decode(decoder TagDecoder, outputStruct interface{}) error {
	v := reflect.ValueOf(outputStruct)
	t := v.Type()
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("expected struct pointer but got: %T", outputStruct)
	}
	t = t.Elem()

	info, err := getStructInfo(t)
	if err != nil {
		return err
	}

	for i := 0; i < t.NumField(); i++ {
		rawValue, err := decoder.DecodeField(info[i])
		if err != nil {
			return fmt.Errorf("error decoding field %v: %s", t.Field(i), err)
		}

		v.Elem().Field(i).Set(reflect.ValueOf(rawValue))
	}
	return nil
}

var structInfoCache = map[reflect.Type][]Field{}

func getStructInfo(t reflect.Type) ([]Field, error) {
	info, found := structInfoCache[t]
	if found {
		return info, nil
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("can only get struct info from structs, but got: %v", t)
	}

	info = []Field{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}

		info = append(info, Field{
			Tags: tags.ParseTags(string(field.Tag)),
			Name: field.Name,
			Type: field.Type,
			Kind: field.Type.Kind(),
		})
	}

	structInfoCache[t] = info
	return info, nil
}
