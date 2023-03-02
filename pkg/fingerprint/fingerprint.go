package fingerprint

import (
	"errors"

	"github.com/Escape-Technologies/goctopus/internal/config"
	out "github.com/Escape-Technologies/goctopus/pkg/output"
)

func FingerprintUrl(url string, fp Fingerprinter, config *config.Config) (*out.FingerprintOutput, error) {
	isGraphql := fp.Graphql()
	if isGraphql {
		if config.Introspection {
			hasIntrospection := fp.Introspection()
			// if hasIntrospection {
			// 	hasFieldSuggestion := fp.FieldSuggestion()

			// }
			return &out.FingerprintOutput{
				Url:           url,
				Introspection: hasIntrospection,
				Type:          out.ResultIsGraphql,
			}, nil
		}
		return &out.FingerprintOutput{
			Url:  url,
			Type: out.ResultIsGraphql,
		}, nil
	}
	return nil, errors.New("no graphql endpoint found on this route")
}

// @todo FingerprintMaybeGraphql
