package goctopus

import (
	"github.com/Escape-Technologies/goctopus/pkg/config"
	out "github.com/Escape-Technologies/goctopus/pkg/output"
	log "github.com/sirupsen/logrus"
)

func FingerprintFromSlice(addresses []string) {
	maxWorkers := config.Get().MaxWorkers
	output := make(chan *out.FingerprintOutput, maxWorkers)
	addressChan := make(chan string, maxWorkers)

	go func() {
		for _, address := range addresses {
			addressChan <- address
		}
		close(addressChan)
	}()

	go FingerprintAddresses(addressChan, output)

	foundCount := out.HandleOutput(output)
	log.Infof("Done. Found %d graphql endpoints", foundCount)
}
