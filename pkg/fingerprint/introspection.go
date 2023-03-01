package fingerprint

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var (
	IntrospectionPayload = []byte(`{"query": "query { __schema { queryType { name } } }"}`)
)

func IsValidIntrospectionResponse(resp *fasthttp.Response) bool {
	if resp.StatusCode() != 200 {
		log.Debugf("Recived status code %v: %v", resp.StatusCode())
		return false
	}
	body := resp.Body()

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
	if err := json.Unmarshal(body, &result); err != nil {
		return false
	}
	if result.Data.Schema.QueryType.Name == "" {
		return false
	}
	return true
}
