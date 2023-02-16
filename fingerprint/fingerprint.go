package fingerprint

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

func SendPost(url string, body []byte) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
	req.SetBody(body)
	req.SetTimeout(2 * time.Second)
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := fasthttp.Do(req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

type Data struct {
	Typename string `json:"__typename" validate:"required"`
}

type Response struct {
	Data Data `json:"data" validate:"required"`
}

func fingerprintUrl(url string) bool {
	body := []byte(`{"query":"{__typename}"}`)
	res, err := SendPost(url, body)
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

func FingerprintDomain(baseDomain string) (string, error) {
	routes := []string {
		"",
		"graphql",
		"api/graphql",
		"api/v1/graphql",
		"api/v2/graphql",
		"appsync",
		"altair",
		"graph",
		"graphql/v1",
		"graphql/v2",
		"api/graphql",
	}
	for _, route := range routes {
		uri := fmt.Sprintf("https://%s/%s", baseDomain, route)
		isGraphql := fingerprintUrl(uri)
		if isGraphql {
			fmt.Println(uri)
			return uri, nil
		}
	}
	return "", errors.New("no graphql endpoint found")
}