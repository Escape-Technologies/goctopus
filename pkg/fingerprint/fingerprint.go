package fingerprint

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Escape-Technologies/goctopus/internal/http"
)

type Data struct {
	Typename string `json:"__typename" validate:"required"`
}

type Response struct {
	Data Data `json:"data" validate:"required"`
}

func fingerprintUrl(url string) bool {
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

func FingerprintDomain(baseDomain string) (string, error) {
	routes := []string{
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
			return uri, nil
		}
	}
	return "", errors.New("no graphql endpoint found")
}
