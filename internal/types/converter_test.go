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
			desc:               "should report error if types are not compatible",
			input:              10,
			targetType:         reflect.TypeOf(struct{}{}),
			expectErrToContain: []string{"cannot convert", "int", "10", "struct"},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			v, err := NewConverter(test.input).Convert(test.targetType)
			if test.expectErrToContain != nil {
				tt.AssertErrContains(t, err, test.expectErrToContain...)
				t.Skip()
			}

			tt.AssertEqual(t, v.Interface(), test.expectedOutput)
		})
	}
}

func intPtr(i int) *int {
	return &i
}
