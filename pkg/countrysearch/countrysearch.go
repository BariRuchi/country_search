package countrysearch

import (
	"CountrySearch/inbound"
	"CountrySearch/lib/cache"
	"CountrySearch/logs"
	"CountrySearch/pkg/apicacheclient"
	"CountrySearch/pkg/apicacheclient/helper"
	"context"
	"encoding/json"
	"fmt"
)

type CountrySearch struct {
	inbound        inbound.CountrySearchInput
	lruCacheClient *cache.LRUCache
}

func (cs *CountrySearch) ServeRequest(ctx context.Context) (string, bool) {

	logs.LogAccess(fmt.Sprintf("Fetching Response For Country : %s", cs.inbound.Name))

	var valid = true
	apiClient := helper.New()

	apiCacheClient := apicacheclient.New(apiClient, cs.lruCacheClient, cs.inbound)
	countryData, err := apiCacheClient.GetCountryData(ctx)
	if err != nil {
		countryData.Error = err.Error()
		valid = false
	}

	finalResponse, err := json.Marshal(countryData)
	if err != nil {
		logs.LogError(fmt.Errorf("error unmarshaling response:%s", err.Error()))
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
