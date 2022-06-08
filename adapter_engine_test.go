package tagmapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vingarcia/tagmapper"
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
}
