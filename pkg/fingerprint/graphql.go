package fingerprint

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/Escape-Technologies/goctopus/internal/http"
)

var (
	GraphqlPayload = []byte(`{"query":"{__typename}"}`)
)

func (fp *fingerprinter) OpenGraphql() (bool, error) {
	body := &GraphqlPayload
	res, err := fp.Client.Post(fp.url, *body)
	if err != nil {
		log.Debugf("Error from %v: %v", fp.url, err)
		return false, err
	}
	log.Debugf("Response from %v: %v", fp.url, res.StatusCode)
	return isOpenGraphqlResponse(res), nil
}

func (fp *fingerprinter) AuthentifiedGraphql() (bool, error) {
	body := &GraphqlPayload
	res, err := fp.Client.Post(fp.url, *body)
	if err != nil {
		log.Debugf("Error from %v: %v", fp.url, err)
		return false, err
	}
	log.Debugf("Response from %v: %v", fp.url, res.StatusCode)
	return isAuthentifiedGraphqlResponse(res), nil
}

func IsOpenGraphql(resp *http.Response) bool {
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

func isOpenGraphqlResponse(resp *http.Response) bool {
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

func isAuthentifiedGraphqlResponse(resp *http.Response) bool {
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
