package fingerprint

import (
	"errors"

	"github.com/Escape-Technologies/goctopus/internal/config"
)

func FingerprintUrl(url string, fp Fingerprinter, config *config.Config) (*FingerprintOutput, error) {
	isGraphql := fp.Graphql()
	if isGraphql {
		if config.Introspection {
			hasIntrospection := fp.Introspection()
			return &FingerprintOutput{
				Url:           url,
				Introspection: hasIntrospection,
				Type:          ResultIsGraphql,
			}, nil
		}
		return &FingerprintOutput{
			Url:  url,
			Type: ResultIsGraphql,
		}, nil
	}
	return nil, errors.New("no graphql endpoint found on this route")
}

// @todo FingerprintMaybeGraphql
