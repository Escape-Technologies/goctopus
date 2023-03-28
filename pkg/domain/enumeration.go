package domain

import (
	"io"

	"github.com/Escape-Technologies/goctopus/pkg/address"
	"github.com/Escape-Technologies/goctopus/pkg/config"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	log "github.com/sirupsen/logrus"
)

func makeCallback(domain *address.Addr, subDomains chan *address.Addr) func(s *resolve.HostEntry) {
	return func(s *resolve.HostEntry) {
		subDomains <- address.NewSourced(s.Host, domain.Source)
	}
}

var (
	runnerInstance *runner.Runner
)

func EnumerateSubdomains(domain *address.Addr, subDomains chan *address.Addr, threads int) (err error) {
	log.Errorf("Enumerating subdomains for %s with %d threads", domain, threads)
	subDomains <- domain
	c := config.Get()

	if !c.SubdomainEnumeration {
		return nil
	}

	if runnerInstance == nil {
		runnerInstance, _ = runner.NewRunner(&runner.Options{
			Threads:            threads,                  // Thread controls the number of threads to use for active enumerations
			Timeout:            c.Timeout,                // Timeout is the seconds to wait for sources to respond
			MaxEnumerationTime: 5,                        // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
			Resolvers:          resolve.DefaultResolvers, // Use the default list of resolvers by marshaling it to the config
			ResultCallback:     makeCallback(domain, subDomains),
			Silent:             true,
		})
	}

	err = runnerInstance.EnumerateSingleDomain(domain.Address, []io.Writer{})
	return err
}
