package structscanner

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/vingarcia/structscanner/internal/types"
	"github.com/vingarcia/structscanner/tags"
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
	idx int

	Tags map[string]string
	Name string
	Kind reflect.Kind
	Type reflect.Type

	IsEmbeded bool
}

func Decode(outputStruct interface{}, decoder TagDecoder) error {
	v := reflect.ValueOf(outputStruct)
	t := v.Type()
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("expected struct pointer but got: %T", outputStruct)
	}
	t = t.Elem()

	fields, err := getStructInfo(t)
	if err != nil {
		return err
	}

	for _, field := range fields {
		rawValue, err := decoder.DecodeField(field)
		if err != nil {
			return fmt.Errorf("error decoding field %v: %s", t.Field(field.idx), err)
		}

		if rawValue == nil {
			continue
		}

		decoder, ok := rawValue.(TagDecoder)
		if ok {
			err := Decode(v.Elem().Field(field.idx).Addr().Interface(), decoder)
			if err != nil {
				return err
			}
			continue
		}

		convertedValue, err := types.NewConverter(rawValue).Convert(field.Type)
		if err != nil {
			return err
		}

		v.Elem().Field(field.idx).Set(convertedValue)
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
		// If it is unexported:
		if unicode.IsLower(rune(field.Name[0])) {
			continue
		}

		parsedTags, err := tags.ParseTags(field.Tag)
		if err != nil {
			return nil, err
		}

		info = append(info, Field{
			idx:  i,
			Tags: parsedTags,
			Name: field.Name,
			Type: field.Type,
			Kind: field.Type.Kind(),

			// ("Anonymous" is the name for embeded fields on the stdlib)
			IsEmbeded: field.Anonymous,
		})
	}

	structInfoCache[t] = info
	return info, nil
}
