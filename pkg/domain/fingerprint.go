package domain

import (
	"errors"
	"net"

	"github.com/Escape-Technologies/goctopus/pkg/address"
	"github.com/Escape-Technologies/goctopus/pkg/endpoint"
	"github.com/Escape-Technologies/goctopus/pkg/output"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// @todo test this
func FingerprintSubDomain(domain *address.Sourced) (*output.FingerprintOutput, error) {
	endpoints := endpoint.FuzzRoutes(domain)

	for _, url := range endpoints {
		output, err := endpoint.FingerprintEndpoint(url)
		// @todo close client
		// fp.Close()

		// At the first timeout, drop the domain
		// @todo number of tries in the config
		if err != nil {
			// If the domain is not a graphql endpoint, continue
			if errors.Is(err, endpoint.ErrNotGraphql) {
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
		output.Domain = domain.Address
		output.Source = domain.Source
		return output, nil
	}
	return nil, errors.New("no graphql endpoint found")
}
