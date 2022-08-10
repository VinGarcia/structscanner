# Intro

This project was created to make it easy to write code that
scans data into structs in safe and efficient manner.

So to make it clear, this is not a library like:

- https://github.com/mitchellh/mapstructure

Nor something like:

https://github.com/spf13/viper

This is a library for allowing you to write your own Viper
or Mapstructure libraries with ease and in a few lines of code,
so that you get exactly what you need and in the way you need it.

So the examples below are examples of things you can get by using
this library. Both examples are also public so you can use them
directly if you want.

But the interesting part is that both were written
in very few lines of code, so check that out too.

## Usage Examples:

The code below will fill the struct with data from env variables.

It will use the `env` tags to map which env var should be used
as source for each of the attributes of the struct.

```golang
// This one is stateless and can be reused:
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
// This one has state and maps a single map to a struct,
// so you might need to declare a new decoder for each use:
var user struct {
	ID       int    `map:"id"`
	Username string `map:"username"`
	Address  struct {
		Street  string `map:"street"`
		City    string `map:"city"`
		Country string `map:"country"`
	} `map:"address"`
}
err := structscanner.Decode(&user, structscanner.NewMapTagDecoder("map", map[string]interface{}{
	"id":       42,
	"username": "fakeUsername",
	"address": map[string]interface{}{
		"street":  "fakeStreet",
		"city":    "fakeCity",
		"country": "fakeCountry",
	},
}))
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
		// it recursively on this nested map:
		return NewMapTagDecoder(e.tagName, nestedMap), nil
	}

	return e.sourceMap[key], nil
}
```

## License

This project was put into public domain, which means you can copy, use and modify
any part of it without mentioning its original source so feel free to do that
if it would be more convenient that way.

Enjoy.
