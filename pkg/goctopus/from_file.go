package goctopus

import (
	"bufio"
	"os"

	"github.com/Escape-Technologies/goctopus/pkg/address"
	"github.com/Escape-Technologies/goctopus/pkg/config"
	out "github.com/Escape-Technologies/goctopus/pkg/output"

	log "github.com/sirupsen/logrus"
)

func FingerprintFromFile(input *os.File) {
	inputBuffer := bufio.NewScanner(input)
	// Scan the file line by line
	inputBuffer.Split(bufio.ScanLines)

	maxWorkers := config.Get().MaxWorkers
	log.Infof("Starting %d workers\n", maxWorkers)
	addresses := make(chan *address.Addr, maxWorkers)
	output := make(chan *out.FingerprintOutput, maxWorkers)

	go FingerprintAddresses(addresses, output)

	go func() {
		for inputBuffer.Scan() {
			addr := inputBuffer.Text()
			addresses <- address.New(addr)
		}
		close(addresses)
	}()

	foundCount := out.HandleOutput(output)
	log.Infof("Done. Found %d graphql endpoints", foundCount)
}
