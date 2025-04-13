package model

type Response struct {
	Name       string `json:"name,omitempty"`
	Capital    string `json:"capital,omitempty"`
	Currency   string `json:"currency,omitempty"`
	Population int    `json:"population,omitempty"`
	Error      string `json:"error,omitempty"`
}
