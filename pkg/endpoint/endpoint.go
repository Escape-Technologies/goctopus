package endpoint

import (
	"github.com/Escape-Technologies/goctopus/pkg/address"
	"github.com/Escape-Technologies/goctopus/pkg/graphql"
	"github.com/Escape-Technologies/goctopus/pkg/http"
	"github.com/Escape-Technologies/goctopus/pkg/introspection"
	"github.com/Escape-Technologies/goctopus/pkg/suggestion"
)

type _endpointFingerprinter struct {
	url    *address.Addr
	client http.Client
}

type endpointFingerprinter interface {
	IsOpenGraphql() (bool, error)
	IsAuthentifiedGraphql() (bool, error)
	HasFieldSuggestion() (bool, error)
	HasIntrospectionOpen() (bool, error)
}

func NewEndpointFingerprinter(url *address.Addr, client http.Client) endpointFingerprinter {
	return &_endpointFingerprinter{
		url:    url,
		client: client,
	}
}

func (e *_endpointFingerprinter) IsOpenGraphql() (bool, error) {
	return graphql.FingerprintOpenGraphql(e.url.Address, e.client)
}

func (e *_endpointFingerprinter) IsAuthentifiedGraphql() (bool, error) {
	return graphql.FingerprintAuthentifiedGraphql(e.url.Address, e.client)
}

func (e *_endpointFingerprinter) HasFieldSuggestion() (bool, error) {
	return suggestion.FingerprintFieldSuggestion(e.url.Address, e.client), nil
}

func (e *_endpointFingerprinter) HasIntrospectionOpen() (bool, error) {
	return introspection.FingerprintIntrospection(e.url.Address, e.client)
}
