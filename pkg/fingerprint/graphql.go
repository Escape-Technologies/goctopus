package fingerprint

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

var (
	GraphqlPayload = []byte(`{"query":"{__typename}"}`)
)

func IsValidGraphqlResponse(resp *fasthttp.Response) bool {
	if resp.StatusCode() != 200 {
		return false
	}
	body := resp.Body()

	type Response struct {
		Data struct {
			Typename string `json:"__typename"`
		} `json:"data"`
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return false
	}
	if result.Data.Typename == "" {
		return false
	}
	return true
}
