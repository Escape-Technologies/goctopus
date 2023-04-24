package goctopus

import (
	"github.com/Escape-Technologies/goctopus/pkg/address"
	out "github.com/Escape-Technologies/goctopus/pkg/output"
)

// Fingerprints a single address and outputs a slice of FingerprintOutput.
func FingerprintAddress(addr *address.Addr) []*out.FingerprintOutput {
	outputSlice := make([]*out.FingerprintOutput, 0)
	output := make(chan *out.FingerprintOutput, 1)
	addresses := make(chan *address.Addr, 1)
	addresses <- addr
	close(addresses)
	go FingerprintAddresses(addresses, output)
	for fingerprintOutput := range output {
		outputSlice = append(outputSlice, fingerprintOutput)
	}
	return outputSlice
}
