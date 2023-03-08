package crawl

import (
	"errors"
	"fmt"
	"net"

	"github.com/Escape-Technologies/goctopus/internal/config"
	"github.com/Escape-Technologies/goctopus/pkg/fingerprint"
	out "github.com/Escape-Technologies/goctopus/pkg/output"
	log "github.com/sirupsen/logrus"

	"github.com/valyala/fasthttp"
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
		"/api",
	}
	// @todo refactor this
	for _, route := range routes {
		url := fmt.Sprintf("https://%s/%s", domain, route)
		fp := fingerprint.NewFingerprinter(url, domain)
		output, err := fingerprint.FingerprintUrl(url, fp, config.Conf)

		// At the first timeout, drop the domain
		// @todo number of tries in the config
		if err != nil {
			// If the domain is not a graphql endpoint, continue
			if errors.Is(err, fingerprint.ErrNotGraphql) {
				continue
			}

			// At the first timeout, drop the domain
			// @todo number of tries in the config
			if errors.Is(err, fasthttp.ErrTimeout) {
				log.Infof("Timeout on %s, skipping.", domain)
				return nil, err
			}

			// If the host can't be resolved, drop the domain
			var dnsErr *net.DNSError
			if errors.As(err, &dnsErr) {
				log.Infof("DNSError on %s, skipping.", domain)
				return nil, err
			}

			// Unknown error
			log.Warnf("Unhandled error on %s, skipping. %v", domain, err)
			return nil, err
		}
		if err == nil {
			output.Domain = domain
			return output, nil
		}
	}
	return nil, errors.New("no graphql endpoint found")
}
