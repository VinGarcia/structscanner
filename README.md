# Intro

This project was created to make it easy to write code that
scans data into structs in safe and efficient manner.

The most important feature of this library is that the user
never needs to deal directly with any of the edge concepts of
the reflect package, for example there is no risk of panics happening
and you only have access to two types from the reflect library:

- reflect.Kind
- reflect.Type

Which you should use only for checking if the input type is what you expect or not.

## Usage Examples:

The code below will fill the struct with data from env variables.

It will use the `env` tags to map which env var should be used
as source for each of the attributes of the struct.

```golang
	decoder := structscanner.FuncTagDecoder(func(field structscanner.Field) (interface{}, error) {
		return os.Getenv(field.Tags["env"]), nil
	})

	var config struct {
		GoPath string `env:"GOPATH"`
		Path   string `env:"PATH"`
		Home   string `env:"HOME"`
	}
	err := structscanner.Decode(&config, decoder)
```

The above example loads data from a global state into the struct.

This second example will fill a struct with the values of an input map:

```golang
	decoder := structscanner.NewMapTagDecoder("map", map[string]interface{}{
		"id":       42,
		"username": "fakeUsername",
		"address": map[string]interface{}{
			"street":  "fakeStreet",
			"city":    "fakeCity",
			"country": "fakeCountry",
		},
	})

	var user struct {
		ID       int    `map:"id"`
		Username string `map:"username"`
		Address  struct {
			Street  string `map:"street"`
			City    string `map:"city"`
			Country string `map:"country"`
		} `map:"address"`
	}
	err := structscanner.Decode(&user, decoder)
```

The code for `FuncTagDecoder` and `MapTagDecoder` are very simple and are also good examples
of how to use this library if you want something slightly different than the examples above:

```golang
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
```

## TODO

- Add test checking if pointers to nested structs work as they should
