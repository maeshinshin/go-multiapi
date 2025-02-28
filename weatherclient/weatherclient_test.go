package weatherclient

import (
	"errors"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	_ "github.com/joho/godotenv/autoload"
)

func TestFetchWeatherData(t *testing.T) {
	type constraint struct {
		noAPIKey       bool
		failtoParsing  bool
		failtoFetching bool
	}

	tests := []struct {
		name       string
		constraint *constraint
		city       string
		err        error
	}{
		{
			name: "not set apiKey",
			constraint: &constraint{
				noAPIKey: true,
			},
			city: "",
			err:  newApiKeyNotFoundError(),
		},
		{
			name: "empty city name",
			city: "",
			err:  newCityParameterNotFoundError(),
		},
		{
			name: "fail to parse API URL",
			constraint: &constraint{
				failtoParsing: true,
			},
			city: "Tokyo",
			err:  newParsingAPIURLFailedError(&url.Error{}),
		},
		{
			name: "fail to fetch weather data",
			constraint: &constraint{
				failtoFetching: true,
			},
			city: "Tokyo",
			err:  newFetchingWeatherDataFailedError(&url.Error{}),
		},
		{
			name: "success fetch weather data",
			city: "Tokyo",
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// var apiKey string
			if tt.constraint != nil {
				if tt.constraint.noAPIKey {
					apiKey := os.Getenv("OPENWEATHER_API_KEY")
					os.Setenv("OPENWEATHER_API_KEY", "")
					defer func() {
						os.Setenv("OPENWEATHER_API_KEY", apiKey)
					}()
				}
				if tt.constraint.failtoParsing || tt.constraint.failtoFetching {
					var tmp string
					apiKey := os.Getenv("OPENWEATHER_API_KEY")
					os.Setenv("OPENWEATHER_API_KEY", "testapi")

					if tt.constraint.failtoFetching {
						apiURL, tmp = "https://gotest.maesh.dev?q=%s&appid=%s", apiURL
					} else {
						apiURL, tmp = "api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", apiURL
					}

					defer func() {
						apiURL = tmp
						os.Setenv("OPENWEATHER_API_KEY", apiKey)
					}()
				}
			}

			_, err := FetchWeatherData(tt.city)
			if diff := cmp.Diff(tt.err, err, cmp.Comparer(func(_, _ error) bool {
				switch e := err.(type) {
				case *FetchingWeatherDataFailedError:
					_, ok := e.err.(*url.Error)
					return ok
				case *ParsingAPIURLFailedError:
					_, ok := e.err.(*url.Error)
					return ok
				default:
					return errors.Is(err, tt.err)
				}
			})); diff != "" {
				t.Errorf("Test %q failed (-want +got):\n%s", tt.name, diff)
			}
		})
	}
}
