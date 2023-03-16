package introspection

import (
	"encoding/json"

	"github.com/Escape-Technologies/goctopus/pkg/http"
	log "github.com/sirupsen/logrus"
)

var (
	IntrospectionPayload = []byte(`{"query": "query { __schema { queryType { name } } }"}`)
)

func FingerprintIntrospection(url string, client http.Client) (bool, error) {
	body := &IntrospectionPayload
	res, err := client.Post(url, *body)
	if err != nil {
		log.Debugf("Error from %v: %v", url, err)
		return false, err
	}
	log.Debugf("Response from %v: %v", url, res.StatusCode)
	return IsValidIntrospectionResponse(res), nil
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
