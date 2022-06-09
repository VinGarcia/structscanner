package types

import (
	"fmt"
	"reflect"
)

// Converter was created to make it easier
// to handle conversion between ptr and non ptr types, e.g.:
//
// - *type to *type
// - type to *type
// - *type to type
// - type to type
type Converter struct {
	BaseType  reflect.Type
	BaseValue reflect.Value
	ElemType  reflect.Type
	ElemValue reflect.Value
}

// NewConverter instantiates a Converter from
// an empty interface.
//
// The input argument can be of any type, but
// if it is a pointer then its Elem() will be
// used as source value for the Converter.Convert()
// method.
func NewConverter(v interface{}) Converter {
	if v == nil {
		// This is necessary so that reflect.ValueOf
		// returns a valid reflect.Value
		v = (*interface{})(nil)
	}

	baseValue := reflect.ValueOf(v)
	baseType := reflect.TypeOf(v)

	elemType := baseType
	elemValue := baseValue
	if baseType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
		elemValue = elemValue.Elem()
	}
	return Converter{
		BaseType:  baseType,
		BaseValue: baseValue,
		ElemType:  elemType,
		ElemValue: elemValue,
	}
}

// Convert attempts to convert the ElemValue to the destType received
// as argument and then returns the converted reflect.Value or an error
func (p Converter) Convert(destType reflect.Type) (reflect.Value, error) {
	destElemType := destType
	if destType.Kind() == reflect.Ptr {
		destElemType = destType.Elem()
	}

	// Return 0 valued destType instance:
	if p.BaseType.Kind() == reflect.Ptr && p.BaseValue.IsNil() {
		// Note that if destType is a ptr it will return a nil ptr.
		return reflect.New(destType).Elem(), nil
	}

	if !p.ElemType.ConvertibleTo(destElemType) {
		return reflect.Value{}, fmt.Errorf(
			"cannot convert from type %v to type %v", p.BaseType, destType,
		)
	}

	destValue := p.ElemValue.Convert(destElemType)

	// Get the address of destValue if necessary:
	if destType.Kind() == reflect.Ptr {
		if !destValue.CanAddr() {
			tmp := reflect.New(destElemType)
			tmp.Elem().Set(destValue)
			destValue = tmp
		} else {
			destValue = destValue.Addr()
		}
	}

	return destValue, nil
}
