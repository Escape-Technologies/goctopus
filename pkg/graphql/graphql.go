package graphql

import (
	"encoding/json"

	"github.com/Escape-Technologies/goctopus/pkg/http"
	log "github.com/sirupsen/logrus"
)

var (
	GraphqlPayload = []byte(`{"query":"{__typename}"}`)
)

func IsOpenGraphqlResponse(resp *http.Response) bool {
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

func FingerprintOpenGraphql(url string, client http.Client) (bool, error) {
	body := &GraphqlPayload
	res, err := client.Post(url, *body)
	if err != nil {
		log.Debugf("Error from %v: %v", url, err)
		return false, err
	}
	log.Debugf("Response from %v: %v", url, res.StatusCode)
	return IsOpenGraphqlResponse(res), nil
}

func IsAuthentifiedGraphqlResponse(resp *http.Response) bool {
	body := resp.Body

	type Response struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	var result Response
	if err := json.Unmarshal(*body, &result); err != nil {
		return false
	}

	if len(result.Errors) < 1 {
		return false
	}

	if result.Errors[0].Message == "" {
		return false
	}

	return true
}

func FingerprintAuthentifiedGraphql(url string, client http.Client) (bool, error) {
	body := &GraphqlPayload
	res, err := client.Post(url, *body)
	if err != nil {
		log.Debugf("Error from %v: %v", url, err)
		return false, err
	}
	log.Debugf("Response from %v: %v", url, res.StatusCode)
	return IsAuthentifiedGraphqlResponse(res), nil
}
