package workers

import (
	"bufio"
	"sync"

	"github.com/Escape-Technologies/goctopus/pkg/crawl"
	out "github.com/Escape-Technologies/goctopus/pkg/output"

	log "github.com/sirupsen/logrus"
)

func worker(subDomains chan string, output chan *out.FingerprintOutput, workerId int, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debugf("Worker %d instantiated\n", workerId)
	for domain := range subDomains {
		log.Debugf("Worker %d started on: %v\n", workerId, domain)
		// res, err := crawl.CrawlSubDomain(domain)
		res, err := crawl.CrawlSubDomain(domain)
		if err == nil {
			log.Debugf("Worker %d found endpoint: %v\n", workerId, res)
			output <- res
		}
	}
	log.Debugf("Worker %d finished\n", workerId)
}

func Orchestrator(inputBuffer *bufio.Scanner, maxWorkers int, output chan *out.FingerprintOutput, count int) {

	domains := make(chan string, 1)
	subDomains := make(chan string, maxWorkers)
	wg := sync.WaitGroup{}
	wg.Add(maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		go worker(subDomains, output, i, &wg)
	}

	i := 1
	for inputBuffer.Scan() {
		domain := inputBuffer.Text()
		log.Infof("(%d/%d) Adding %v to the queue\n", i, count, domain)
		subDomains <- domain
		crawl.CrawlDomain(domain, subDomains, true)
		i++
	}

	close(domains)
	log.Debugf("Orchestrator finished, waiting for workers to finish...")
	wg.Wait()
	close(output)
	log.Debugf("All workers finished")
}
