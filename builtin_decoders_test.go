package structscanner_test

import (
	"os"
	"testing"

	"github.com/vingarcia/structscanner"
	tt "github.com/vingarcia/structscanner/helpers/testtools"
)

func TestFuncTagDecoder(t *testing.T) {
	os.Setenv("GOPATH", "fakeGOPATH")
	os.Setenv("PATH", "fakePATH")
	os.Setenv("HOME", "fakeHOME")

	decoder := structscanner.FuncTagDecoder(func(field structscanner.Field) (interface{}, error) {
		return os.Getenv(field.Tags["env"]), nil
	})

	var config struct {
		GoPath string `env:"GOPATH"`
		Path   string `env:"PATH"`
		Home   string `env:"HOME"`
	}
	err := structscanner.Decode(decoder, &config)
	tt.AssertNoErr(t, err)
	tt.AssertEqual(t, config.GoPath, "fakeGOPATH")
	tt.AssertEqual(t, config.Path, "fakePATH")
	tt.AssertEqual(t, config.Home, "fakeHOME")
}

func TestMapTagDecoder(t *testing.T) {
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
	err := structscanner.Decode(decoder, &user)
	tt.AssertNoErr(t, err)
	tt.AssertEqual(t, user.ID, 42)
	tt.AssertEqual(t, user.Username, "fakeUsername")
	tt.AssertEqual(t, user.Address.Street, "fakeStreet")
	tt.AssertEqual(t, user.Address.City, "fakeCity")
	tt.AssertEqual(t, user.Address.Country, "fakeCountry")
}
