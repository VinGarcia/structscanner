package tagmapper_test

import (
	"testing"

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
		tt.AssertNoErr(t, err)
		tt.AssertEqual(t, output.Attr1, "fake-value-for-string")
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
		tt.AssertNoErr(t, err)
		tt.AssertEqual(t, output.Attr1, "fake-value-for-string")
		tt.AssertEqual(t, output.Attr2, "placeholder")
	})

	t.Run("should be able to fill multiple attributes", func(t *testing.T) {
		decoder := FuncTagDecoder(func(field tagmapper.Field) (interface{}, error) {
			v := map[string]string{
				"f1": "v1",
				"f2": "v2",
				"f3": "v3",
			}[field.Tags["map"]]

			return v, nil
		})

		var output struct {
			Attr1 string `map:"f1"`
			Attr2 string `map:"f2"`
			Attr3 string `map:"f3"`
		}
		err := tagmapper.Decode(decoder, &output)
		tt.AssertNoErr(t, err)
		tt.AssertEqual(t, output.Attr1, "v1")
		tt.AssertEqual(t, output.Attr2, "v2")
		tt.AssertEqual(t, output.Attr3, "v3")
	})

	t.Run("should ignore private fields", func(t *testing.T) {
		decoder := FuncTagDecoder(func(field tagmapper.Field) (interface{}, error) {
			return 64, nil
		})

		var output struct {
			attr1 int `env:"attr1"`
		}
		err := tagmapper.Decode(decoder, &output)
		tt.AssertNoErr(t, err)
		tt.AssertEqual(t, output.attr1, 0)
	})

	t.Run("should convert types correctly", func(t *testing.T) {
		t.Run("should convert different types of integers", func(t *testing.T) {
			decoder := FuncTagDecoder(func(field tagmapper.Field) (interface{}, error) {
				return uint64(10), nil
			})

			var output struct {
				Attr1 int `env:"attr1"`
			}
			err := tagmapper.Decode(decoder, &output)
			tt.AssertNoErr(t, err)
			tt.AssertEqual(t, output.Attr1, 10)
		})

		t.Run("should convert from ptr to non ptr", func(t *testing.T) {
			decoder := FuncTagDecoder(func(field tagmapper.Field) (interface{}, error) {
				i := 64
				return &i, nil
			})

			var output struct {
				Attr1 int `env:"attr1"`
			}
			err := tagmapper.Decode(decoder, &output)
			tt.AssertNoErr(t, err)
			tt.AssertEqual(t, output.Attr1, 64)
		})

		t.Run("should convert from ptr to non ptr", func(t *testing.T) {
			decoder := FuncTagDecoder(func(field tagmapper.Field) (interface{}, error) {
				return 64, nil
			})

			var output struct {
				Attr1 *int `env:"attr1"`
			}
			err := tagmapper.Decode(decoder, &output)
			tt.AssertNoErr(t, err)
			tt.AssertNotEqual(t, output.Attr1, nil)
			tt.AssertEqual(t, *output.Attr1, 64)
		})

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
