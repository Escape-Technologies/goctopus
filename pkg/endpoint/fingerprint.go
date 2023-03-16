package endpoint

import (
	"errors"

	"github.com/Escape-Technologies/goctopus/pkg/config"
	"github.com/Escape-Technologies/goctopus/pkg/http"
	"github.com/Escape-Technologies/goctopus/pkg/output"
)

var ErrNotGraphql = errors.New("no graphql endpoint found on this route")

func fingerprintEndpoint(url string, e endpointFingerprinter, config *config.Config) (*output.FingerprintOutput, error) {
	out := &output.FingerprintOutput{
		Url:  url,
		Type: output.ResultOpenGraphql,
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

func FingerprintEndpoint(url string) (*output.FingerprintOutput, error) {
	c := config.Get()
	client := http.NewClient(c)
	e := NewEndpointFingerprinter(url, client)
	res, err := fingerprintEndpoint(url, e, c)
	client.DeleteUrlCache(url)
	return res, err
}
