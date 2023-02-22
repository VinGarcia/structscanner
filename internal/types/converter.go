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

	destValue, err := p.convert(destElemType, destType)
	if err != nil {
		return reflect.Value{}, err
	}

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

func (p Converter) convert(destElemType reflect.Type, destType reflect.Type) (reflect.Value, error) {
	if p.ElemType.Kind() == reflect.Map &&
		destElemType.Kind() == reflect.Map &&
		p.ElemType != destElemType {

		return p.convertMap(destElemType, destType)
	}

	if !p.ElemType.ConvertibleTo(destElemType) {
		return reflect.Value{}, fmt.Errorf(
			"cannot convert from type %v to type %v, received value was: %v",
			p.BaseType, destType, p.ElemValue,
		)
	}

	return p.ElemValue.Convert(destElemType), nil
}

func (p Converter) convertMap(destElemType reflect.Type, destType reflect.Type) (reflect.Value, error) {
	destElemKeyType := destElemType.Key()
	destElemValueType := destElemType.Elem()

	targetMap := reflect.MakeMap(destElemType)
	iter := p.ElemValue.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		if key.Type().Kind() == reflect.Interface && !key.IsNil() {
			key = key.Elem()
		}

		if value.Type().Kind() == reflect.Interface && !value.IsNil() {
			value = value.Elem()
		}

		if !key.Type().ConvertibleTo(destElemKeyType) {
			return reflect.Value{}, fmt.Errorf(
				"cannot convert map key '%v' of type %v to target map key of type: %v",
				key, key.Type(), destElemKeyType,
			)
		}

		if !value.Type().ConvertibleTo(destElemValueType) {
			return reflect.Value{}, fmt.Errorf(
				"cannot convert map value: '%v' of type: %v, on key: '%v', to type: %v",
				value, value.Type(), key, destElemValueType,
			)
		}

		targetMap.SetMapIndex(key.Convert(destElemKeyType), value.Convert(destElemValueType))
	}

	return targetMap, nil
}
