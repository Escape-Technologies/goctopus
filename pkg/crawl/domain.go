package crawl

import (
	"io"

	"github.com/Escape-Technologies/goctopus/internal/config"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

func CrawlDomain(domain string, subDomains chan string, enumerationEnabled bool) (err error) {
	subDomains <- domain

	if !enumerationEnabled {
		return nil
	}

	runnerInstance, _ := runner.NewRunner(&runner.Options{
		Threads:            10,                       // Thread controls the number of threads to use for active enumerations
		Timeout:            config.Conf.Timeout,      // Timeout is the seconds to wait for sources to respond
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
