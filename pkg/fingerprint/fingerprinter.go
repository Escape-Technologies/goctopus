package fingerprint

import (
	"github.com/Escape-Technologies/goctopus/internal/http"
)

type Fingerprinter interface {
	OpenGraphql() (bool, error)
	AuthentifiedGraphql() (bool, error)
	IntrospectionOpen() (bool, error)
	FieldSuggestionEnabled() (bool, error)
}

type fingerprinter struct {
	url    string
	Client http.Client
}

func NewFingerprinter(url string) *fingerprinter {
	client := http.NewClient()
	return &fingerprinter{
		url:    url,
		Client: client,
	}
}
