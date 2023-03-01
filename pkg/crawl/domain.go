package crawl

import (
	"errors"
	"fmt"

	"github.com/Escape-Technologies/goctopus/internal/config"
	"github.com/Escape-Technologies/goctopus/pkg/fingerprint"
)

// type DomainCrawler interface {

// func Domain(domain string) (Output, error) {
// 	// @todo add subdomain enumeration here
// 	return FingerprintSubDomain(domain)
// }

// extract those

func CrawlSubDomain(domain string) (fingerprint.FingerprintOutput, error) {
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
	fp := fingerprint.NewUrlFingerprinter(domain)
	for _, route := range routes {
		url := fmt.Sprintf("https://%s/%s", domain, route)
		output, err := CrawlRoute(url, fp)
		if err == nil {
			return output, nil
		}
	}
	return nil, errors.New("no graphql endpoint found")
}

func CrawlRoute(url string, fp fingerprint.UrlFingerprinter) (fingerprint.FingerprintOutput, error) {
	isGraphql := fp.Graphql()
	if isGraphql {
		if config.Conf.Introspection {
			hasIntrospection := fp.Introspection()
			return &fingerprint.IsGraphqlOutput{
				Url:           url,
				Introspection: hasIntrospection,
			}, nil
		}
		return &fingerprint.IsGraphqlOutput{
			Url:           url,
			Introspection: false,
		}, nil
	}
	return nil, errors.New("no graphql endpoint found on this route")
}
