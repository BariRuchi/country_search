package apicall

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ApiCall struct {
	url string
	ctx context.Context
}

func (a *ApiCall) Call() (*http.Response, error) {
	req, err := http.NewRequestWithContext(a.ctx, http.MethodGet, a.url, nil)
	if err != nil {
		return &http.Response{}, fmt.Errorf("request creation error: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return resp, fmt.Errorf("http request error: %w", err)
	}
	return resp, nil
}

func New(ctx context.Context, url string) *ApiCall {
	api := new(ApiCall)
	api.ctx = ctx
	api.url = url
	return api
}
