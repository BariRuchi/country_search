package apicacheclient

import (
	"CountrySearch/apicacheclient/apiclient"
	"CountrySearch/apicacheclient/cache"
	"CountrySearch/inbound"
	"CountrySearch/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type ApiCacheClient struct {
	apiClient *apiclient.ApiClient
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
			err = fmt.Errorf("err while unmarshaling cache reponse : %s", err.Error())
			log.Print(err.Error())
			return response, err
		}
		fmt.Printf("Cache Response : %s", value)
		return response, nil
	}

	countryDetails, err := cli.apiClient.FetchCountryDataFromAPI(ctx, cli.inbound.Name)
	if err != nil {
		err = fmt.Errorf("error while fetching contry data : %s", err.Error())
		log.Print(err.Error())
		return response, err
	}

	cacheVal, err := json.Marshal(countryDetails)
	if err == nil {
		cli.lruCache.Set(key, string(cacheVal))
	}
	return countryDetails, nil
}

func (cli *ApiCacheClient) getCacheKey() string {
	key := strings.ToLower(cli.inbound.Name)
	log.Printf("Cache key :%s", key)
	return key
}

func New(apiClient *apiclient.ApiClient, lruCache *cache.LRUCache, inbound inbound.CountrySearchInput) *ApiCacheClient {
	cli := new(ApiCacheClient)
	cli.apiClient = apiClient
	cli.lruCache = lruCache
	cli.inbound = inbound
	return cli
}
