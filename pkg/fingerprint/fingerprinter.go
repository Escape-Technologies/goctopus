package fingerprint

import (
	"github.com/Escape-Technologies/goctopus/internal/http"
)

type Fingerprinter interface {
	Graphql() (bool, error)
	Introspection() (bool, error)
	FieldSuggestion() (bool, error)
}

type fingerprinter struct {
	url    string
	domain string
	Client http.Client
}

func NewFingerprinter(url string, domain string) *fingerprinter {
	client := http.NewClient()
	return &fingerprinter{
		url:    url,
		domain: domain,
		Client: client,
	}
}
