package domain

import (
	"io"

	"github.com/Escape-Technologies/goctopus/pkg/address"
	"github.com/Escape-Technologies/goctopus/pkg/config"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

func makeCallback(domain *address.Sourced, subDomains chan *address.Sourced) func(s *resolve.HostEntry) {
	return func(s *resolve.HostEntry) {
		subDomains <- address.NewSourced(s.Host, domain.Source)
	}
}

func EnumerateSubdomains(domain *address.Sourced, subDomains chan *address.Sourced) (err error) {
	subDomains <- domain
	c := config.Get()

	if !c.SubdomainEnumeration {
		return nil
	}

	runnerInstance, _ := runner.NewRunner(&runner.Options{
		Threads:            c.MaxWorkers,             // Thread controls the number of threads to use for active enumerations
		Timeout:            c.Timeout,                // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 5,                        // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
		Resolvers:          resolve.DefaultResolvers, // Use the default list of resolvers by marshaling it to the config
		ResultCallback:     makeCallback(domain, subDomains),
		Silent:             true,
	})

	err = runnerInstance.EnumerateSingleDomain(domain.Address, []io.Writer{})
	return err
}
