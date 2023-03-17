package goctopus

import (
	"github.com/Escape-Technologies/goctopus/pkg/address"
	"github.com/Escape-Technologies/goctopus/pkg/config"
	out "github.com/Escape-Technologies/goctopus/pkg/output"
	log "github.com/sirupsen/logrus"
)

func FingerprintFromSlice(addresses []string) {
	maxWorkers := config.Get().MaxWorkers
	output := make(chan *out.FingerprintOutput, maxWorkers)
	addressChan := make(chan *address.Sourced, maxWorkers)

	go func() {
		for _, addr := range addresses {
			addressChan <- address.NewSourced(addr, addr)
		}
		close(addressChan)
	}()

	go FingerprintAddresses(addressChan, output)

	foundCount := out.HandleOutput(output)
	log.Infof("Done. Found %d graphql endpoints", foundCount)
}
