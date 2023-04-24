package address

import "strings"

type Addr struct {
	Address  string
	Source   string
	Metadata map[string]string
}

func New(address string) *Addr {
	return &Addr{
		Address:  address,
		Source:   address,
		Metadata: map[string]string{},
	}
}

// This allows to attach a custom source to an address
func NewSourced(address, source string) *Addr {
	return &Addr{
		Address: address,
		Source:  source,
	}
}

func (a *Addr) AddMetadata(key, value string) {
	a.Metadata[key] = value
}

func (a *Addr) Copy() *Addr {
	metadataCopy := make(map[string]string)
	for k, v := range a.Metadata {
		metadataCopy[k] = v
	}
	return &Addr{
		Address:  strings.Clone(a.Address),
		Source:   strings.Clone(a.Source),
		Metadata: metadataCopy,
	}
}
