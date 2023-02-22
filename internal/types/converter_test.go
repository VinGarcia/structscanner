package types

import (
	"reflect"
	"testing"

	tt "github.com/vingarcia/structscanner/internal/testtools"
)

func TestConverter(t *testing.T) {
	tests := []struct {
		desc               string
		input              any
		targetType         reflect.Type
		expectedOutput     any
		expectErrToContain []string
	}{
		{
			desc:           "should convert int to int",
			input:          10,
			targetType:     reflect.TypeOf(10),
			expectedOutput: 10,
		},
		{
			desc:           "should convert compatible types",
			input:          10.0,
			targetType:     reflect.TypeOf(10),
			expectedOutput: 10,
		},
		{
			desc:           "should convert ptr into types",
			input:          intPtr(10),
			targetType:     reflect.TypeOf(10),
			expectedOutput: 10,
		},
		{
			desc:           "should convert type into ptrs",
			input:          10,
			targetType:     reflect.TypeOf(new(int)),
			expectedOutput: intPtr(10),
		},
		{
			desc:           "should convert ptr into ptr",
			input:          intPtr(10),
			targetType:     reflect.TypeOf(new(int)),
			expectedOutput: intPtr(10),
		},
		{
			desc:           "should return nil for nil inputs if target type is a pointer",
			input:          nil,
			targetType:     reflect.TypeOf(new(int)),
			expectedOutput: (*int)(nil),
		},
		{
			desc:           "should return zero value for nil inputs if target type is not a pointer",
			input:          nil,
			targetType:     reflect.TypeOf(10),
			expectedOutput: 0,
		},
		{
			desc: "should convert maps of different but compatible types",
			input: map[string]string{
				"fakeKey": "fakeValue",
			},
			targetType: reflect.TypeOf(map[any]any{}),
			expectedOutput: map[any]any{
				"fakeKey": "fakeValue",
			},
		},
		{
			desc: "should convert maps with interface subtypes in a best-effort basis", // (Try to convert it, if not possible return an error)
			input: map[any]any{
				"fakeKey": "fakeValue",
			},
			targetType: reflect.TypeOf(map[string]string{}),
			expectedOutput: map[string]string{
				"fakeKey": "fakeValue",
			},
		},
		{
			desc:               "should report error if types are not compatible",
			input:              10,
			targetType:         reflect.TypeOf(struct{}{}),
			expectErrToContain: []string{"cannot convert", "int", "10", "struct"},
		},
		{
			desc: "should report error if map key is not compatible with target map key",
			input: map[struct{}]any{
				struct{}{}: "fakeValue",
			},
			targetType:         reflect.TypeOf(map[string]string{}),
			expectErrToContain: []string{"key", "struct", "string"},
		},
		{
			desc: "should report error if map value is not compatible with target map value",
			input: map[string]struct{}{
				"fakeKey": struct{}{},
			},
			targetType:         reflect.TypeOf(map[string]string{}),
			expectErrToContain: []string{"value", "struct", "string"},
		},
		{
			desc: "should not panic if a map key is nil",
			input: map[any]string{
				nil: "fakeValue",
			},
			targetType:         reflect.TypeOf(map[string]string{}),
			expectErrToContain: []string{"cannot convert", "nil", "to", "string"},
		},
		{
			desc: "should not panic if a map value is nil",
			input: map[string]any{
				"fakeKey": nil,
			},
			targetType:         reflect.TypeOf(map[string]string{}),
			expectErrToContain: []string{"cannot convert", "nil", "to", "string"},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			v, err := NewConverter(test.input).Convert(test.targetType)
			if test.expectErrToContain != nil {
				tt.AssertErrContains(t, err, test.expectErrToContain...)
				t.Skip()
			}

			tt.AssertNoErr(t, err)
			tt.AssertEqual(t, v.Interface(), test.expectedOutput)
		})
	}
}

func intPtr(i int) *int {
	return &i
}
