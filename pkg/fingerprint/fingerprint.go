package fingerprint

import (
	"encoding/json"

	"github.com/Escape-Technologies/goctopus/internal/http"
)

type UrlFingerprinter interface {
	// @todo http client here for mocking
	Graphql() bool
	Introspection() bool
}

type urlFingerprinter struct {
	url string
}

func NewUrlFingerprinter(url string) *urlFingerprinter {
	return &urlFingerprinter{
		url: url,
	}
}

func (fp *urlFingerprinter) Graphql() bool {
	type Response struct {
		Data struct {
			Typename string `json:"__typename"`
		} `json:"data"`
	}

	body := []byte(`{"query":"{__typename}"}`)
	res, err := http.Post(fp.url, body)
	if err != nil {
		return false
	}
	var result Response
	err = json.Unmarshal(res, &result)
	if err != nil {
		return false
	}
	if result.Data.Typename != "" {
		return true
	}
	return false
}

func (fp *urlFingerprinter) Introspection() bool {
	type Response struct {
		Data struct {
			Schema struct {
				QueryType struct {
					Name string `json:"name"`
				} `json:"queryType"`
			} `json:"__schema"`
		} `json:"data"`
	}

	body := []byte(`{"query": "query { __schema { queryType { name } } }"}`)
	res, err := http.Post(fp.url, body)
	if err != nil {
		return false
	}
	var result Response
	if err = json.Unmarshal(res, &result); err != nil {
		return false
	}
	if result.Data.Schema.QueryType.Name != "" {
		return true
	}
	return false
}
