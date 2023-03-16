package goctopus

import (
	"sync"

	"github.com/Escape-Technologies/goctopus/internal/utils"
	"github.com/Escape-Technologies/goctopus/pkg/config"
	"github.com/Escape-Technologies/goctopus/pkg/domain"
	"github.com/Escape-Technologies/goctopus/pkg/endpoint"
	"github.com/Escape-Technologies/goctopus/pkg/output"
	log "github.com/sirupsen/logrus"
)

func worker(addresses chan string, output chan *output.FingerprintOutput, workerId int, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debugf("Worker %d instantiated", workerId)
	for address := range addresses {
		log.Debugf("Worker %d started on: %v", workerId, address)
		res, err := FingerprintAddress(address)
		if err == nil {
			log.Debugf("Worker %d found endpoint: %v", workerId, res)
			output <- res
		}
	}
	log.Debugf("Worker %d finished", workerId)
}

func FingerprintAddress(address string) (*output.FingerprintOutput, error) {
	// If the domain is a url, we don't need to crawl it
	if utils.IsUrl(address) {
		return endpoint.FingerprintEndpoint(address)
	} else {
		return domain.FingerprintSubDomain(address)
	}
}

// An addresses can be a domain or an url
func FingerprintAddresses(addresses chan string, output chan *output.FingerprintOutput) {

	maxWorkers := config.Get().MaxWorkers
	enumeratedAddresses := make(chan string, config.Get().MaxWorkers)

	workersWg := sync.WaitGroup{}
	workersWg.Add(maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		go worker(enumeratedAddresses, output, i, &workersWg)
	}

	i := 1
	for address := range addresses {
		log.Debugf("(%d) Adding %v to the queue", i, address)
		// If the domain is a url, we don't need to crawl it
		if utils.IsUrl(address) {
			enumeratedAddresses <- address
		} else {
			if err := domain.EnumerateSubdomains(address, enumeratedAddresses); err != nil {
				log.Errorf("Error enumerating subdomains for %v: %v", address, err)
			}
		}
		i++
	}

	close(enumeratedAddresses)
	log.Debugf("Waiting for workers to finish...")
	workersWg.Wait()
	close(output)
	log.Debugf("All workers finished")
}
