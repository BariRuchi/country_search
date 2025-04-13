package apiclient

import (
	"CountrySearch/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

var URL = "https://restcountries.com/v3.1/name/"

type ApiClient struct {
}

func (a *ApiClient) FetchCountryDataFromAPI(ctx context.Context, countryName string) (model.Response, error) {
	response := model.Response{}
	url := URL + countryName

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return response, fmt.Errorf("request creation error: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("Oops!! status code: %d", resp.StatusCode)
	}

	var data []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil || len(data) == 0 {
		return response, fmt.Errorf("error decoding response: %w", err)
	}

	countryData := data[0]

	response.Name = countryData["name"].(map[string]interface{})["common"].(string)
	response.Capital = countryData["capital"].([]interface{})[0].(string)
	currencyCode := reflect.ValueOf(countryData["currencies"]).MapKeys()[0].String()
	response.Currency = countryData["currencies"].(map[string]interface{})[currencyCode].(map[string]interface{})["symbol"].(string)
	response.Population = int(countryData["population"].(float64))

	return response, nil

}

func New() *ApiClient {
	return new(ApiClient)
}
