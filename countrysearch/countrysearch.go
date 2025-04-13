package countrysearch

import (
	"CountrySearch/apicacheclient"
	"CountrySearch/apicacheclient/apiclient"
	"CountrySearch/apicacheclient/cache"
	"CountrySearch/inbound"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type CountrySearch struct {
	inbound        inbound.CountrySearchInput
	lruCacheClient *cache.LRUCache
}

func (cs *CountrySearch) ServeRequest(ctx context.Context) (string, bool) {

	var valid = true
	apiClient := apiclient.New()

	apiCacheClient := apicacheclient.New(apiClient, cs.lruCacheClient, cs.inbound)
	countryData, err := apiCacheClient.GetCountryData(ctx)
	if err != nil {
		countryData.Error = err.Error()
		valid = false
	}

	finalResponse, err := json.Marshal(countryData)
	if err != nil {
		err = fmt.Errorf("error unmarshaling response:%s", err.Error())
		log.Print(err.Error())
		valid = false
	}

	return string(finalResponse), valid
}

func New(inbound inbound.CountrySearchInput, lruCacheClient *cache.LRUCache) *CountrySearch {
	cs := new(CountrySearch)
	cs.inbound = inbound
	cs.lruCacheClient = lruCacheClient
	return cs
}
