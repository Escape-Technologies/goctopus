package workers

import (
	"bufio"
	"sync"

	"github.com/Escape-Technologies/goctopus/pkg/crawl"
	"github.com/Escape-Technologies/goctopus/pkg/fingerprint"

	log "github.com/sirupsen/logrus"
)

func worker(domains chan string, output chan *fingerprint.FingerprintOutput, workerId int, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debugf("Worker %d instantiated\n", workerId)
	for domain := range domains {
		log.Debugf("Worker %d started on: %v\n", workerId, domain)
		res, err := crawl.CrawlSubDomain(domain)
		if err == nil {
			log.Debugf("Worker %d found endpoint: %v\n", workerId, res)
			output <- res
		}
	}
	log.Debugf("Worker %d finished\n", workerId)
}

func Orchestrator(inputBuffer *bufio.Scanner, maxWorkers int, output chan *fingerprint.FingerprintOutput, count int) {

	domains := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		go worker(domains, output, i, &wg)
	}

	i := 0
	for inputBuffer.Scan() {
		domain := inputBuffer.Text()
		log.Infof("(%d/%d) Adding %v to the queue\n", i, count, domain)
		domains <- domain
		i++
	}

	close(domains)
	log.Debugf("Orchestrator finished, waiting for workers to finish...")
	wg.Wait()
	close(output)
	log.Debugf("All workers finished")
}
