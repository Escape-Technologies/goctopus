package crawl

import (
	"errors"
	"fmt"

	"github.com/Escape-Technologies/goctopus/internal/config"
	"github.com/Escape-Technologies/goctopus/pkg/fingerprint"
)

func CrawlSubDomain(domain string) (*fingerprint.FingerprintOutput, error) {
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
	// @todo refactor this
	for _, route := range routes {
		url := fmt.Sprintf("https://%s/%s", domain, route)
		fp := fingerprint.NewFingerprinter(url)
		output, err := CrawlRoute(url, fp)
		if err == nil {
			output.Domain = domain
			return output, nil
		}
	}
	return nil, errors.New("no graphql endpoint found")
}

// @note this is the problematic function, hard to test and really important
func CrawlRoute(url string, fp fingerprint.Fingerprinter) (*fingerprint.FingerprintOutput, error) {
	isGraphql := fp.Graphql()
	if isGraphql {
		if config.Conf.Introspection {
			hasIntrospection := fp.Introspection()
			return &fingerprint.FingerprintOutput{
				Url:           url,
				Introspection: hasIntrospection,
				Type:          fingerprint.IsGraphql,
			}, nil
		}
		return &fingerprint.FingerprintOutput{
			Url:  url,
			Type: fingerprint.IsGraphql,
		}, nil
	}
	return nil, errors.New("no graphql endpoint found on this route")
}
