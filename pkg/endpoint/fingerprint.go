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
		Url:          url.Address,
		Source:       url.Source,
		SchemaStatus: output.SchemaStatusClosed,
	}

	isOpenGraphql, err := e.IsOpenGraphql()
	if err != nil {
		return nil, err
	}

	if config.EngineFingerprinting {
		out.Engine = e.GetEngine()
	}

	if !isOpenGraphql {
		isAuthenticatedGraphql, err := e.IsAuthenticatedGraphql()
		if err != nil {
			return nil, err
		}

		if !isAuthenticatedGraphql {
			return nil, ErrNotGraphql
		}

		out.Authenticated = true
		return out, nil
	}

	if config.Introspection {
		hasIntrospectionOpen, err := e.HasIntrospectionOpen()
		if err != nil {
			return nil, err
		}

		if hasIntrospectionOpen {
			out.SchemaStatus = output.SchemaStatusOpen
		} else {
			out.SchemaStatus = output.SchemaStatusClosed
		}
	}

	if out.SchemaStatus == output.SchemaStatusClosed && config.FieldSuggestion {
		hasFieldSuggestion, err := e.HasFieldSuggestion()
		if err != nil {
			return nil, err
		}

		if hasFieldSuggestion {
			out.SchemaStatus = output.SchemaStatusLeaking
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
