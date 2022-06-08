package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Foo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var err error
var pathParam int
var headerParam string
var body Foo

var weight = 10

type FuncTagDecoder func(info Field) (interface{}, error)

func (e FuncTagDecoder) DecodeField(info Field) (interface{}, error) {
	return e(info)
}

func TestUnmarshal(t *testing.T) {
	t.Run("should parse a single tag with a hardcoded value", func(t *testing.T) {
		decoder := FuncTagDecoder(func(field Field) (interface{}, error) {
			return "fake-value-for-string", nil
		})

		var output struct {
			Attr1 string `env:"attr1"`
		}
		err := Decode(decoder, &output)
		assert.Nil(t, err)
		assert.Equal(t, output.Attr1, "fake-value-for-string")
	})
}
