package workers

import (
	"bufio"
	"goctopus/pkg/fingerprint"
	"sync"

	log "github.com/sirupsen/logrus"
)

func worker(domains chan string, endpoints chan string, workerId int, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debugf("Worker %d instantiated\n", workerId)
	for domain := range domains {
		log.Debugf("Worker %d started on: %v\n", workerId, domain)
		endpoint, err := fingerprint.FingerprintDomain(domain)
		if err == nil {
			log.Debugf("Worker %d found endpoint: %v\n", workerId, endpoint)
			endpoints <- endpoint
		}	
	}
	log.Debugf("Worker %d finished\n", workerId)
}

func Orchestrator(inputBuffer *bufio.Scanner, maxWorkers int, endpoints chan string, count int) {

	domains := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		go worker(domains, endpoints, i, &wg)
	}

	i := 0
	for inputBuffer.Scan() {
		domain := inputBuffer.Text()
		log.Infof("(%d/%d) Adding %v to the queue\n", i,count, domain)
		domains <- domain
		i++
	}

	close(domains)
	log.Debugf("Orchestrator finished, waiting for workers to finish...")
	wg.Wait()
	close(endpoints)
	log.Debugf("All workers finished")
}