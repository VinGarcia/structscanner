[![CI](https://github.com/VinGarcia/structscanner/actions/workflows/ci.yml/badge.svg)](https://github.com/VinGarcia/structscanner/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/VinGarcia/structscanner/branch/master/graph/badge.svg?token=5CNJ867C66)](https://codecov.io/gh/VinGarcia/structscanner)
[![Go Reference](https://pkg.go.dev/badge/github.com/vingarcia/structscanner.svg)](https://pkg.go.dev/github.com/vingarcia/structscanner)
![Go Report Card](https://goreportcard.com/badge/github.com/vingarcia/structscanner)

# Intro

This project was created to make it easy to write code that
scans data into structs in a safe and efficient manner.

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

## Understanding the Project:

![image showing that the TagDecoder interface is a wrapper that
adapts a data source so that the Decode function can use it to fill
the attributes of a struct](docs/understanding-the-project.png)

So we have 3 pieces here:

The data source can be anything in your context: A source map,
env variables, a config file, and so on, you name it.

By default the `structscanner` library does not know how to interact
with the data source that you have chosen, so you have to teach it.

That's where the decoder comes in:
This decoder should be a wrapper over your chosen data source,
and it should implement the `structscanner.TagDecoder` interface,
so that when requested it will read the data source on behalf
of the `structcanner` library.

> Note: It will probably be necessary to instantiate a new Decoder for each
> instance of a data source, which I know feels least than ideal but
> it was necessary for fully decoupling from the data source.
>
> It is also easy enough to write a function that does the instantiation of the wrapper,
> and then calls the `Decode()` function, like in the examples below so it looks
> better for the final user.

Having your decoder instantiated you can now call the `structscanner.Decode()`
function passing the decoder instance and the target struct that you want
to be filled with data, and the `Decode()` function will handle all the
necessary reflection magic for you.

It will also keep a cache with the most expensive steps (the ones that use reflection the most)
so that decoding can be done efficiently.

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
// so you might need to instantiate a new decoder for each input map:
var user struct {
	ID       int    `map:"id"`
	Username string `map:"username"`
	Address  struct {
		Street  string `map:"street"`
		City    string `map:"city"`
		Country string `map:"country"`
	} `map:"address"`
    SomeSlice []int `map:"some_slice"`
}
err := structscanner.Decode(&user, structscanner.NewMapTagDecoder("map", map[string]interface{}{
	"id":       42,
	"username": "fakeUsername",
	"address": map[string]interface{}{
		"street":  "fakeStreet",
		"city":    "fakeCity",
		"country": "fakeCountry",
	},
    // Note that even though the type of the slice below
    // differs from the struct slice it will convert all
    // values correctly:
    "some_slice": []float64{1.0, 2.0, 3.0},
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

If you wish to use the Field info (names, tags, type etc) elsewhere you can use the `GetStructInfo()` function.

```golang

type User struct {
	Name    string `map:"name"`
	HomeDir string `map:"home"`
}

info, err := structscanner.GetStructInfo(&User{})
if err != nil {
	panic(err)
}

for _, field := range info.Fields {
	fmt.Println("Field %q has tags %v", field.Name, field.Tags)
}
```

It is possible to pass a `reflection.Type` object to `GetStructInfo`, which is particularly useful for nested structs:

```golang

type Address struct {}
type User struct {
	Name    string `map:"name"`
	HomeDir Address `map:"home"`
}

info, err := structscanner.GetStructInfo(&User{})
if err != nil {
	panic(err)
}

for _, field := range info.Fields {
	fmt.Println("Field %q has tags %v", field.Name, field.Tags)
	if field.Kind == reflect.Struct {
		nestedInfo, err := structscanner.GetStructInfo(field.Type)
		fmt.Println("Nested Field %q has %d fields", field.Name, len(nestedInfo.Fields))
	}
}
```


## License

This project was put into public domain, which means you can copy, use and modify
any part of it without mentioning its original source so feel free to do that
if it would be more convenient that way.

Enjoy.
