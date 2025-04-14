package apicacheclient

import (
	"CountrySearch/inbound"
	"CountrySearch/lib/cache"
	"CountrySearch/logs"
	"CountrySearch/model"
	"CountrySearch/pkg/apicacheclient/helper"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type ApiCacheClient struct {
	apiClient *helper.ApiClient
	lruCache  *cache.LRUCache
	inbound   inbound.CountrySearchInput
}

func (cli *ApiCacheClient) GetCountryData(ctx context.Context) (model.Response, error) {
	var response model.Response

	key := cli.getCacheKey()

	if val, ok := cli.lruCache.Get(key); ok {
		value := val.(string)
		err := json.Unmarshal([]byte(value), &response)
		if err != nil {
			logs.LogError(fmt.Errorf("err while unmarshaling cache reponse : %s", err.Error()))
			return response, err
		}

		logs.LogAccess(fmt.Sprintf("Cache Response : %s", value))
		return response, nil
	}

	countryDetails, err := cli.apiClient.FetchCountryDataFromAPI(ctx, cli.inbound.Name)
	if err != nil {
		logs.LogError(fmt.Errorf("error while fetching contry data : %s", err.Error()))
		return response, err
	}

	cacheVal, err := json.Marshal(countryDetails)
	if err != nil {
		logs.LogError(fmt.Errorf("error while unmarshaling Response : %s", err.Error()))
		return response, err
	}

	logs.LogAccess(fmt.Sprintf("Api Call Response : %s", cacheVal))

	val := string(cacheVal)
	logs.LogAccess(fmt.Sprintf("Setting Response In cache : %s", val))
	cli.lruCache.Set(key, val)

	return countryDetails, nil
}

func (cli *ApiCacheClient) getCacheKey() string {
	key := strings.ToLower(cli.inbound.Name)
	logs.LogAccess(fmt.Sprintf("Cache Key : %s", key))
	return key
}

func New(apiClient *helper.ApiClient, lruCache *cache.LRUCache, inbound inbound.CountrySearchInput) *ApiCacheClient {
	cli := new(ApiCacheClient)
	cli.apiClient = apiClient
	cli.lruCache = lruCache
	cli.inbound = inbound
	return cli
}
