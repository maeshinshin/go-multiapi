package weatherclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var apiURL = "https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric"

func FetchWeatherData(city string) (*WeatherData, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")

	if apiKey == "" {
		return nil, newApiKeyNotFoundError()
	}

	if city == "" {
		return nil, newCityParameterNotFoundError()
	}

	uRL, err := url.ParseRequestURI(fmt.Sprintf(apiURL, city, apiKey))
	if err != nil {
		return nil, newParsingAPIURLFailedError(err)
	}

	resp, err := http.Get(uRL.String())
	if err != nil {
		return nil, newFetchingWeatherDataFailedError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, newApiRequestFailedError(resp.StatusCode)
	}

	var weatherData WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, err
	}

	return &weatherData, nil
}
