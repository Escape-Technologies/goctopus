package fingerprint

import (
	"encoding/json"

	"github.com/Escape-Technologies/goctopus/internal/http"

	log "github.com/sirupsen/logrus"
)

var (
	IntrospectionPayload = []byte(`{"query": "query { __schema { queryType { name } } }"}`)
)

func (fp *fingerprinter) Introspection() bool {
	body := &IntrospectionPayload
	res, err := fp.Client.Post(fp.url, *body)
	log.Debugf("Response from %v: %v", fp.url, string(*res.Body))
	if err != nil {
		log.Debugf("Error from %v: %v", fp.url, err)
		return false
	}
	return IsValidIntrospectionResponse(res)
}

func IsValidIntrospectionResponse(resp *http.Response) bool {
	if resp.StatusCode != 200 {
		return false
	}

	type Response struct {
		Data struct {
			Schema struct {
				QueryType struct {
					Name string `json:"name"`
				} `json:"queryType"`
			} `json:"__schema"`
		} `json:"data"`
	}

	var result Response
	if err := json.Unmarshal(*resp.Body, &result); err != nil {
		return false
	}
	if result.Data.Schema.QueryType.Name == "" {
		return false
	}
	return true
}
