package workers

import (
	"bufio"
	"sync"

	"github.com/Escape-Technologies/goctopus/internal/config"
	"github.com/Escape-Technologies/goctopus/internal/utils"
	"github.com/Escape-Technologies/goctopus/pkg/crawl"
	"github.com/Escape-Technologies/goctopus/pkg/fingerprint"
	out "github.com/Escape-Technologies/goctopus/pkg/output"

	log "github.com/sirupsen/logrus"
)

// @todo refactor this
func worker(addresses chan string, output chan *out.FingerprintOutput, workerId int, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debugf("Worker %d instantiated", workerId)
	for address := range addresses {
		log.Debugf("Worker %d started on: %v", workerId, address)
		var (
			res *out.FingerprintOutput
			err error
		)
		if utils.IsUrl(address) {
			log.Debugf("Worker %d found url: %v", workerId, address)
			fp := fingerprint.NewFingerprinter(address)
			res, err = fingerprint.FingerprintUrl(address, fp, config.Conf)
			fp.Close()
		} else {
			res, err = crawl.CrawlSubDomain(address)
		}
		if err == nil {
			log.Debugf("Worker %d found endpoint: %v", workerId, res)
			output <- res
		}
	}
	log.Debugf("Worker %d finished", workerId)
}

func Orchestrator(inputBuffer *bufio.Scanner, maxWorkers int, output chan *out.FingerprintOutput, count int) {

	// Adresses can be subdomains or urls
	addresses := make(chan string, maxWorkers)
	workersWg := sync.WaitGroup{}
	workersWg.Add(maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		go worker(addresses, output, i, &workersWg)
	}

	i := 1
	for inputBuffer.Scan() {
		domain := inputBuffer.Text()
		log.Debugf("(%d/%d) Adding %v to the queue", i, count, domain)
		// If the domain is a url, we don't need to crawl it
		if utils.IsUrl(domain) {
			addresses <- domain
		} else {
			crawl.CrawlDomain(domain, addresses, config.Conf)
		}
		i++
	}

	close(addresses)
	log.Debugf("Orchestrator finished, waiting for workers to finish...")
	workersWg.Wait()
	close(output)
	log.Debugf("All workers finished")
}
