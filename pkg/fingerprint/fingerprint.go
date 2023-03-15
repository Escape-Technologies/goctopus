package fingerprint

import (
	"errors"

	"github.com/Escape-Technologies/goctopus/internal/config"
	"github.com/Escape-Technologies/goctopus/pkg/output"
)

var ErrNotGraphql = errors.New("no graphql endpoint found on this route")

func FingerprintUrl(url string, fp Fingerprinter, config *config.Config) (*output.FingerprintOutput, error) {
	out := &output.FingerprintOutput{
		Url:  url,
		Type: output.ResultOpenGraphql,
	}
	isOpenGraphql, err := fp.OpenGraphql()
	if err != nil {
		return nil, err
	}

	if !isOpenGraphql {
		isAuthentifiedGraphql, err := fp.AuthentifiedGraphql()

		if err != nil {
			return nil, err
		}

		if !isAuthentifiedGraphql {
			return nil, ErrNotGraphql
		}

		out.Type = output.ResultAuthentifiedGraphql
	}

	if isOpenGraphql && config.Introspection {
		out.Introspection, err = fp.IntrospectionOpen()
		if err != nil {
			return nil, err
		}
	}

	if out.Introspection && config.FieldSuggestion {
		out.FieldSuggestion, err = fp.FieldSuggestionEnabled()
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// @todo FingerprintMaybeGraphql
