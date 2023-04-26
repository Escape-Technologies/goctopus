package address

import (
	"strings"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

type Addr struct {
	Address     string
	Source      string
	Metadata    map[string]string
	callback    func() // callback to call when the origin is done (copies == 0)
	copiesCount *int32 // number of copies of this address
	done        bool   // true if done has been called, to prevent calling it twice
}

func New(address string) *Addr {
	copiesCount := int32(1)
	addr := &Addr{
		Address:     address,
		Source:      address,
		Metadata:    make(map[string]string),
		copiesCount: &copiesCount,
		done:        false,
	}
	addr.callback = func() {
		log.Infof("Done fingerprinting %s", addr.Address)
	}
	return addr
}

func (a *Addr) SetCallback(callback func(addr *Addr)) {
	a.callback = func() { callback(a) }
}

func (a *Addr) AddMetadata(key, value string) {
	a.Metadata[key] = value
}

func (a *Addr) Copy() *Addr {
	if a.done {
		panic("copy called on a done address")
	}
	metadataCopy := make(map[string]string)
	for k, v := range a.Metadata {
		metadataCopy[k] = v
	}
	atomic.AddInt32(a.copiesCount, 1)
	child := &Addr{
		Address:     strings.Clone(a.Address),
		Source:      strings.Clone(a.Source),
		Metadata:    metadataCopy,
		copiesCount: a.copiesCount,
		callback:    a.callback,
		done:        false,
	}
	return child
}

func (a *Addr) Derive(newAddress string) *Addr {
	addr := a.Copy()
	addr.Address = newAddress
	return addr
}

// don't call done on the children
func (a *Addr) Done() {
	if a.done {
		return
	}
	a.done = true
	copies := atomic.AddInt32(a.copiesCount, -1)
	// @todo remove this
	if copies <= 0 {
		a.callback()
	}
}
