package crawl

import (
	"errors"
	"fmt"

	"github.com/Escape-Technologies/goctopus/internal/config"
	"github.com/Escape-Technologies/goctopus/pkg/fingerprint"
	out "github.com/Escape-Technologies/goctopus/pkg/output"
)

func CrawlSubDomain(domain string) (*out.FingerprintOutput, error) {
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
		fp := fingerprint.NewFingerprinter(url, domain)
		output, err := fingerprint.FingerprintUrl(url, fp, config.Conf)
		if err == nil {
			output.Domain = domain
			return output, nil
		}
	}
	return nil, errors.New("no graphql endpoint found")
}
