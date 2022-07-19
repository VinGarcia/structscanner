package structscanner

import (
	"fmt"
	"reflect"
)

// FuncTagDecoder is a simple wrapper for decoders that do not need
// to keep any state.
type FuncTagDecoder func(info Field) (interface{}, error)

// DecodeField implements the TagDecoder interface
func (e FuncTagDecoder) DecodeField(info Field) (interface{}, error) {
	return e(info)
}

// MapTagDecoder can be used to fill a struct with the values of a map.
//
// It works recursively so you can pass nested structs to it.
type MapTagDecoder struct {
	tagName   string
	sourceMap map[string]any
}

// NewMapTagDecoder returns a new decoder for filling a given struct
// with the values from the sourceMap argument.
//
// The values from the sourceMap will be mapped to the struct using the key
// present in the tagName of each field of the struct.
func NewMapTagDecoder(tagName string, sourceMap map[string]interface{}) MapTagDecoder {
	return MapTagDecoder{
		tagName:   tagName,
		sourceMap: sourceMap,
	}
}

// DecodeField implements the TagDecoder interface
func (e MapTagDecoder) DecodeField(info Field) (interface{}, error) {
	key := info.Tags[e.tagName]
	if info.Kind == reflect.Struct {
		nestedMap, ok := e.sourceMap[key].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(
				"can't map %T into nested struct %s of type %v",
				e.sourceMap[key], info.Name, info.Type,
			)
		}

		// By returning a decoder you tell the library to run
		// it recursively on top of this attribute:
		return NewMapTagDecoder(e.tagName, nestedMap), nil
	}

	return e.sourceMap[key], nil
}
