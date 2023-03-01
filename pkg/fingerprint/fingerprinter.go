package fingerprint

import (
	log "github.com/sirupsen/logrus"

	"github.com/Escape-Technologies/goctopus/internal/http"
)

type Fingerprinter interface {
	Graphql() bool
	Introspection() bool
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

func (fp *fingerprinter) Graphql() bool {
	body := &GraphqlPayload
	res, err := fp.Client.Post(fp.url, *body)
	log.Debugf("Response from %v: %v", fp.url, string(res.Body()))
	if err != nil {
		log.Debugf("Error from %v: %v", fp.url, err)
		return false
	}
	return IsValidGraphqlResponse(res)
}

func (fp *fingerprinter) Introspection() bool {
	body := &IntrospectionPayload
	res, err := fp.Client.Post(fp.url, *body)
	if err != nil {
		return false
	}
	return IsValidIntrospectionResponse(res)
}
