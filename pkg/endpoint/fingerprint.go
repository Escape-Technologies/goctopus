package endpoint

import (
	"errors"

	"github.com/Escape-Technologies/goctopus/pkg/address"
	"github.com/Escape-Technologies/goctopus/pkg/config"
	"github.com/Escape-Technologies/goctopus/pkg/http"
	"github.com/Escape-Technologies/goctopus/pkg/output"
)

var ErrNotGraphql = errors.New("no graphql endpoint found on this route")

func fingerprintEndpoint(url *address.Addr, e endpointFingerprinter, config *config.Config) (*output.FingerprintOutput, error) {
	out := &output.FingerprintOutput{
		Url:    url.Address,
		Source: url.Source,
		Type:   output.ResultOpenGraphql,
	}
	isOpenGraphql, err := e.IsOpenGraphql()
	if err != nil {
		return nil, err
	}

	if !isOpenGraphql {
		isAuthentifiedGraphql, err := e.IsAuthentifiedGraphql()

		if err != nil {
			return nil, err
		}

		if !isAuthentifiedGraphql {
			return nil, ErrNotGraphql
		}

		out.Type = output.ResultAuthentifiedGraphql
	}

	if isOpenGraphql && config.Introspection {
		out.Introspection, err = e.HasIntrospectionOpen()
		if err != nil {
			return nil, err
		}
	}

	if !out.Introspection && config.Introspection && config.FieldSuggestion {
		out.FieldSuggestion, err = e.HasFieldSuggestion()
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func FingerprintEndpoint(addr *address.Addr) (*output.FingerprintOutput, error) {
	c := config.Get()
	client := http.NewClient(c)
	e := NewEndpointFingerprinter(addr, client)
	res, err := fingerprintEndpoint(addr, e, c)
	client.DeleteUrlCache(addr.Address)
	return res, err
}
