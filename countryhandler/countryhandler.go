package countryhandler

import (
	"CountrySearch/apicacheclient/cache"
	"CountrySearch/countrysearch"
	"CountrySearch/inbound"
	"net/http"
	"sync"
	"time"
)

const MaxCapacity = 20
const ttl = 2 * time.Hour

type Handler struct {
	lruCacheClientOnce sync.Once
	lruCacheClient     *cache.LRUCache
}

func (ch *Handler) CountryHandler(w http.ResponseWriter, r *http.Request) {

	inbound := inbound.CountrySearchInput{}
	inbound.Name = r.URL.Query().Get("name")
	if inbound.Name == "" {
		http.Error(w, "Missing 'name' query param", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	app := countrysearch.New(inbound, ch.lruCacheClient)
	appResponse, isValidResponse := app.ServeRequest(ctx)
	if isValidResponse {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	_, _ = w.Write([]byte(appResponse))
	return
}

func (ch *Handler) setCacheClient() *cache.LRUCache {
	ch.lruCacheClientOnce.Do(func() {
		ch.lruCacheClient = cache.NewLRUCache(MaxCapacity, ttl)
	})
	return ch.lruCacheClient
}

func New() *Handler {
	ch := new(Handler)
	ch.setCacheClient()
	return ch
}
