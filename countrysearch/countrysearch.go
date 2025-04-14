package countrysearch

import (
	"CountrySearch/apicacheclient"
	"CountrySearch/apicacheclient/helper"
	"CountrySearch/inbound"
	"CountrySearch/lib/cache"
	"CountrySearch/logs"
	"context"
	"encoding/json"
	"fmt"
)

type CountrySearch struct {
	inbound        inbound.CountrySearchInput
	lruCacheClient *cache.LRUCache
	errorLogChan   chan logs.LogMessage
	accessLogChan  chan logs.LogMessage
}

func (cs *CountrySearch) ServeRequest(ctx context.Context) (string, bool) {

	cs.accessLogChan <- logs.CreateAccessLogMessage(fmt.Sprintf("Fetching Response For Country : %s", cs.inbound.Name))

	var valid = true
	apiClient := helper.New()

	apiCacheClient := apicacheclient.New(apiClient, cs.lruCacheClient, cs.inbound, cs.errorLogChan, cs.accessLogChan)
	countryData, err := apiCacheClient.GetCountryData(ctx)
	if err != nil {
		countryData.Error = err.Error()
		valid = false
	}

	finalResponse, err := json.Marshal(countryData)
	if err != nil {
		cs.errorLogChan <- logs.CreateErrorLogMessage(fmt.Errorf("error unmarshaling response:%s", err.Error()))
		valid = false
	}

	return string(finalResponse), valid
}

func New(inbound inbound.CountrySearchInput, lruCacheClient *cache.LRUCache, errorLogChan, accessLogChan chan logs.LogMessage) *CountrySearch {
	cs := new(CountrySearch)
	cs.inbound = inbound
	cs.lruCacheClient = lruCacheClient
	cs.errorLogChan = errorLogChan
	cs.accessLogChan = accessLogChan
	return cs
}
