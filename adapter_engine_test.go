package tagmapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vingarcia/tagmapper"
	tt "github.com/vingarcia/tagmapper/helpers/testtools"
)

type FuncTagDecoder func(info tagmapper.Field) (interface{}, error)

func (e FuncTagDecoder) DecodeField(info tagmapper.Field) (interface{}, error) {
	return e(info)
}

func TestUnmarshal(t *testing.T) {
	t.Run("should parse a single tag with a hardcoded value", func(t *testing.T) {
		decoder := FuncTagDecoder(func(field tagmapper.Field) (interface{}, error) {
			return "fake-value-for-string", nil
		})

		var output struct {
			Attr1 string `env:"attr1"`
		}
		err := tagmapper.Decode(decoder, &output)
		assert.Nil(t, err)
		assert.Equal(t, output.Attr1, "fake-value-for-string")
	})

	t.Run("should ignore attributes if the function returns a nil value", func(t *testing.T) {
		decoder := FuncTagDecoder(func(field tagmapper.Field) (interface{}, error) {
			envTag := field.Tags["env"]
			if envTag == "" {
				return nil, nil
			}

			return "fake-value-for-string", nil
		})

		var output struct {
			Attr1 string `env:"attr1"`
			Attr2 string `someothertag:"attr2"`
		}
		output.Attr2 = "placeholder"
		err := tagmapper.Decode(decoder, &output)
		assert.Nil(t, err)
		assert.Equal(t, output.Attr1, "fake-value-for-string")
		assert.Equal(t, output.Attr2, "placeholder")
	})

	t.Run("should report errors correctly", func(t *testing.T) {
		tests := []struct {
			desc               string
			value              interface{}
			targetStruct       interface{}
			expectErrToContain []string
		}{
			{
				desc:  "should report error if the type doesnt match",
				value: "example-value",
				targetStruct: &struct {
					Attr1 int `env:"attr1"`
				}{},
				expectErrToContain: []string{"string", "int"},
			},
		}
		for _, test := range tests {
			t.Run(test.desc, func(t *testing.T) {
				decoder := FuncTagDecoder(func(field tagmapper.Field) (interface{}, error) {
					return test.value, nil
				})

				err := tagmapper.Decode(decoder, test.targetStruct)
				tt.AssertErrContains(t, err, test.expectErrToContain...)
			})
		}
	})
}
