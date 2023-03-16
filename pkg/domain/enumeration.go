package domain

import (
	"io"

	"github.com/Escape-Technologies/goctopus/pkg/config"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

func EnumerateSubdomains(domain string, subDomains chan string) (err error) {
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
		ResultCallback: func(s *resolve.HostEntry) { // Callback function to execute for available host
			subDomains <- s.Host
		},
		Silent: true,
	})

	err = runnerInstance.EnumerateSingleDomain(domain, []io.Writer{})
	return err
}
