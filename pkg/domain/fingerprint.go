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
func FingerprintSubDomain(domain *address.Addr) (*output.FingerprintOutput, error) {
	endpoints := endpoint.FuzzRoutes(domain)

	for _, addr := range endpoints {
		output, err := endpoint.FingerprintEndpoint(addr)
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
				log.Debugf("Timeout on %s, skipping.", domain)
			}

			// If the host can't be resolved, drop the domain
			var dnsErr *net.DNSError
			if errors.As(err, &dnsErr) {
				log.Debugf("DNSError on %s, skipping.", domain)
			}

			// Unknown error
			log.Debugf("Unhandled error on %s, not skipping. %v", domain, err)
			return nil, err
		}
		output.Domain = domain.Address
		output.Source = domain.Source
		output.Metadata = domain.Metadata
		return output, nil
	}
	return nil, errors.New("no graphql endpoint found")
}
