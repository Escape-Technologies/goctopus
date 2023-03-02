package fingerprint

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/Escape-Technologies/goctopus/internal/http"
)

var (
	GraphqlPayload = []byte(`{"query":"{__typename}"}`)
)

func (fp *fingerprinter) Graphql() bool {
	body := &GraphqlPayload
	res, err := fp.Client.Post(fp.url, *body)
	if err != nil {
		log.Debugf("Error from %v: %v", fp.url, err)
		return false
	}
	log.Debugf("Response from %v: %v", fp.url, string(*res.Body))
	return IsValidGraphqlResponse(res)
}

func IsValidGraphqlResponse(resp *http.Response) bool {
	if resp.StatusCode != 200 {
		return false
	}
	body := resp.Body

	type Response struct {
		Data struct {
			Typename string `json:"__typename"`
		} `json:"data"`
	}

	var result Response
	if err := json.Unmarshal(*body, &result); err != nil {
		return false
	}
	if result.Data.Typename == "" {
		return false
	}
	return true
}
