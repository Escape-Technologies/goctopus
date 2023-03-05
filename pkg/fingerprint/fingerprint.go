package fingerprint

import (
	"errors"

	"github.com/Escape-Technologies/goctopus/internal/config"
	out "github.com/Escape-Technologies/goctopus/pkg/output"
)

var ErrNotGraphql = errors.New("no graphql endpoint found on this route")

func FingerprintUrl(url string, fp Fingerprinter, config *config.Config) (*out.FingerprintOutput, error) {
	out := &out.FingerprintOutput{
		Url:  url,
		Type: out.ResultIsGraphql,
	}
	isGraphql, err := fp.Graphql()
	if err != nil {
		return nil, err
	}

	if !isGraphql {
		return nil, ErrNotGraphql
	}
	if isGraphql && config.Introspection {
		out.Introspection, err = fp.Introspection()
		if err != nil {
			return nil, err
		}
	}
	if out.Introspection && config.FieldSuggestion {
		out.FieldSuggestion, err = fp.FieldSuggestion()
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// @todo FingerprintMaybeGraphql
