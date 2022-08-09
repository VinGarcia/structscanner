package structscanner_test

import (
	"os"
	"testing"

	"github.com/vingarcia/structscanner"
	tt "github.com/vingarcia/structscanner/internal/testtools"
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
	err := structscanner.Decode(&config, decoder)
	tt.AssertNoErr(t, err)
	tt.AssertEqual(t, config.GoPath, "fakeGOPATH")
	tt.AssertEqual(t, config.Path, "fakePATH")
	tt.AssertEqual(t, config.Home, "fakeHOME")
}

func TestMapTagDecoder(t *testing.T) {
	t.Run("should work for valid structs", func(t *testing.T) {
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
		tt.AssertNoErr(t, err)

		tt.AssertEqual(t, user.ID, 42)
		tt.AssertEqual(t, user.Username, "fakeUsername")
		tt.AssertEqual(t, user.Address.Street, "fakeStreet")
		tt.AssertEqual(t, user.Address.City, "fakeCity")
		tt.AssertEqual(t, user.Address.Country, "fakeCountry")
	})

	t.Run("should return error if we try to save something that is not a map into a nested struct", func(t *testing.T) {
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
			"address":  "notAMap",
		}))

		tt.AssertErrContains(t, err, "string", "Address", "Street", "City", "Country")
	})
}
