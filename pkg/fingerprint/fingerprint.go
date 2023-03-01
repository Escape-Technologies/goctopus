package fingerprint

import (
	"encoding/json"

	"github.com/Escape-Technologies/goctopus/internal/http"
)

// @todo graphql and introspection could be in separate files

type Fingerprinter interface {
	Graphql() bool
	Introspection() bool
}

type fingerprinter struct {
	url    string
	client http.Client
	//@todo http client here for mocking
}

func NewFingerprinter(url string) *fingerprinter {
	client := http.NewClient()
	return &fingerprinter{
		url:    url,
		client: client,
	}
}

func (fp *fingerprinter) Graphql() bool {
	type Response struct {
		Data struct {
			Typename string `json:"__typename"`
		} `json:"data"`
	}

	body := []byte(`{"query":"{__typename}"}`)
	res, err := fp.client.Post(fp.url, body)
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

func (fp *fingerprinter) Introspection() bool {
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
	res, err := fp.client.Post(fp.url, body)
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
