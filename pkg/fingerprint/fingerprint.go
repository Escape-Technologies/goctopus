package fingerprint

import (
	"errors"

	"github.com/Escape-Technologies/goctopus/internal/config"
	out "github.com/Escape-Technologies/goctopus/pkg/output"
)

func FingerprintUrl(url string, fp Fingerprinter, config *config.Config) (*out.FingerprintOutput, error) {
	out := &out.FingerprintOutput{
		Url:  url,
		Type: out.ResultIsGraphql,
	}
	isGraphql := fp.Graphql()
	if !isGraphql {
		return nil, errors.New("no graphql endpoint found on this route")
	}
	if isGraphql && config.Introspection {
		out.Introspection = fp.Introspection()
	}
	if out.Introspection && config.FieldSuggestion {
		out.FieldSuggestion = fp.FieldSuggestion()
	}
	return out, nil
}

// @todo FingerprintMaybeGraphql
