package fingerprint

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Escape-Technologies/goctopus/internal/config"
	"github.com/Escape-Technologies/goctopus/internal/http"
)

func fingerprintGraphql(url string) bool {
	type Response struct {
		Data struct {
			Typename string `json:"__typename"`
		} `json:"data"`
	}

	body := []byte(`{"query":"{__typename}"}`)
	res, err := http.Post(url, body)
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

func fingerprintIntrospection(url string) bool {

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
	res, err := http.Post(url, body)
	if err != nil {
		return false
	}
	var result Response
	err = json.Unmarshal(res, &result)
	if err != nil {
		return false
	}
	if result.Data.Schema.QueryType.Name != "" {
		return true
	}
	return false
}

func FingerprintDomain(domain string) (Output, error) {
	routes := []string{
		"",
		"graphql",
		"api/graphql",
		"api/v2/graphql",
		"api/v1/graphql",
		"appsync",
		"altair",
		"graph",
		"graphql/v2",
		"graphql/v1",
		"api/graphql",
	}
	// @todo add subdomain enumeration here
	// @todo refactor this
	for _, route := range routes {
		url := fmt.Sprintf("https://%s/%s", domain, route)
		isGraphql := fingerprintGraphql(url)
		if isGraphql {
			if config.Conf.Introspection {
				hasIntrospection := fingerprintIntrospection(url)
				return IsGraphqlOutput{
					Url:           url,
					Introspection: hasIntrospection,
					Domain:        domain,
				}, nil
			}
			return IsGraphqlOutput{
				Url:           url,
				Introspection: false,
				Domain:        domain,
			}, nil
		}
	}
	return nil, errors.New("no graphql endpoint found")
}
