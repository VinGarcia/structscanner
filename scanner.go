package structscanner

import (
	"fmt"
	"reflect"
	"sync"
	"unicode"

	"github.com/vingarcia/structscanner/internal/types"
	"github.com/vingarcia/structscanner/tags"
)

// TagDecoder is the interface that allows the Decode function to get values
// from any data source and then use these values to fill a targetStruct.
//
// The struct that implements this TagDecoder interface should
// handle each call to `DecodeField()` by returning the value
// that should be written to the Field described in the `field` argument.
//
// The Decode() function will then take care of checking and making any
// necessary conversions between the returned value and actual struct field.
//
// The FuncTagDecoder and MapTagDecoder are examples of how this interface
// can be implemented, please read the source code of these two types
// to better understand this interface.
type TagDecoder interface {
	DecodeField(field Field) (interface{}, error)
}

// Field is the input expected by the `DecodeField` method
// of the TagDecoder interface and contains all the information
// about the field that is currently being targeted by the
// Decode() function.
type Field struct {
	idx int

	Tags map[string]string
	Name string
	Kind reflect.Kind
	Type reflect.Type

	IsEmbeded bool
}

// Decode reads from the input decoder in order to fill the
// attributes of an target struct.
func Decode(targetStruct interface{}, decoder TagDecoder) error {
	v := reflect.ValueOf(targetStruct)
	t := v.Type()
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("expected struct pointer but got: %T", targetStruct)
	}
	if v.IsNil() {
		return fmt.Errorf("expected non-nil pointer to struct, but got: %#v", targetStruct)
	}

	t = t.Elem()

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("can only get struct info from structs, but got: %#v", targetStruct)
	}

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
			fieldType := v.Elem().Field(field.idx).Type()
			fieldAddr := v.Elem().Field(field.idx).Addr()
			if fieldType.Kind() == reflect.Ptr {
				if fieldAddr.Elem().IsNil() {
					// If this field is a nil pointer, do struct.Field = new(*T):
					fieldAddr.Elem().Set(reflect.New(fieldType.Elem()))
				}
				// Now since it is a pointer, drop one level for the
				// decode function to receive a *struct instead of a **struct:
				fieldAddr = fieldAddr.Elem()
			}

			err := Decode(fieldAddr.Interface(), decoder)
			if err != nil {
				return fmt.Errorf("error decoding nested field '%s': %s", t.Field(field.idx).Name, err)
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

// This cache is kept as a pkg variable
// because the total number of types on a program
// should be finite. So keeping a single cache here
// works fine.
var structInfoCache = &sync.Map{}

func getStructInfo(t reflect.Type) ([]Field, error) {
	data, _ := structInfoCache.Load(t)
	info, ok := data.([]Field)
	if ok {
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

	structInfoCache.Store(t, info)
	return info, nil
}
