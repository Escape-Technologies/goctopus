package goctopus

import (
	"context"
	"sync"

	"github.com/Escape-Technologies/goctopus/internal/utils"
	"github.com/Escape-Technologies/goctopus/pkg/address"
	"github.com/Escape-Technologies/goctopus/pkg/config"
	"github.com/Escape-Technologies/goctopus/pkg/domain"
	"github.com/Escape-Technologies/goctopus/pkg/endpoint"
	"github.com/Escape-Technologies/goctopus/pkg/output"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

func worker(addresses chan *address.Addr, output chan *output.FingerprintOutput, workerId int, sem *semaphore.Weighted) {
	log.Debugf("Worker %d instantiated", workerId)
	for address := range addresses {
		sem.Acquire(context.Background(), 1)
		log.Debugf("Worker %d started on: %v", workerId, address)
		res, err := FingerprintAddress(address)
		if err == nil {
			log.Debugf("Worker %d found endpoint: %v", workerId, res)
			output <- res
		}
		sem.Release(1)
	}
	log.Debugf("Worker %d finished", workerId)
}

func FingerprintAddress(address *address.Addr) (*output.FingerprintOutput, error) {
	// If the domain is a url, we don't need to crawl it
	if utils.IsUrl(address.Address) {
		return endpoint.FingerprintEndpoint(address)
	} else {
		return domain.FingerprintSubDomain(address)
	}
}

func asyncEnumeration(address *address.Addr, enumeratedAddresses chan *address.Addr, threads int, sem *semaphore.Weighted, wg *sync.WaitGroup) {
	defer wg.Done()
	defer sem.Release(int64(threads))
	if err := domain.EnumerateSubdomains(address, enumeratedAddresses, threads); err != nil {
		log.Errorf("Error enumerating subdomains for %v: %v", address, err)
	}
}

// An addresses can be a domain or an url
func FingerprintAddresses(addresses chan *address.Addr, output chan *output.FingerprintOutput) {

	maxWorkers := config.Get().MaxWorkers
	enumeratedAddresses := make(chan *address.Addr, config.Get().MaxWorkers)

	sem := semaphore.NewWeighted(int64(maxWorkers))
	enumerationWg := sync.WaitGroup{}
	enumerationThreads := utils.MinInt(maxWorkers, 10)

	for i := 0; i < maxWorkers; i++ {
		go worker(enumeratedAddresses, output, i, sem)
	}

	i := 1
	for address := range addresses {
		log.Debugf("(%d) Adding %v to the queue", i, address)
		// If the domain is a url, we don't need to crawl it
		if utils.IsUrl(address.Address) {
			sem.Acquire(context.Background(), 1)
			enumeratedAddresses <- address
		} else {
			// 10 threads for subdomain enumeration, unless maxWorkers is less than 10
			enumerationWg.Add(1)
			sem.Acquire(context.Background(), int64(enumerationThreads))
			log.Errorf("%v", address)
			go asyncEnumeration(address, enumeratedAddresses, enumerationThreads, sem, &enumerationWg)
		}
		i++
	}

	enumerationWg.Wait()
	close(enumeratedAddresses)
	log.Debugf("Waiting for workers to finish...")
	sem.Acquire(context.Background(), int64(maxWorkers))
	close(output)
	log.Debugf("All workers finished")
}
