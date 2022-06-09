# Intro

This project was created to make it easy to write code that
scans data into structs in safe and efficient manner.

## Usage Example:

The code below will fill the struct with data from env variables.

It will use the `env` tags to map which env var should be used
as source for each of the attributes of the struct.

```golang
decoder := FuncTagDecoder(func(field structscanner.Field) (interface{}, error) {
	return os.Getenv(field.Tags["env"]]), nil
})

var configs struct {
	GoPath string `env:"GOPATH"`
	Path   string `env:"PATH"`
	Home   string `env:"HOME"`
}
err := structscanner.Decode(decoder, &output)
```
