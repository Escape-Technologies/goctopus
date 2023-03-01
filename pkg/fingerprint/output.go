package fingerprint

import (
	"encoding/json"

	"github.com/Escape-Technologies/goctopus/internal/config"
)

type FingerprintResult string

const (
	ResultIsGraphql    FingerprintResult = "IS_GRAPHQL"
	ResultMaybeGraphql FingerprintResult = "MAYBE_GRAPHQL"
)

type FingerprintOutput struct {
	Type          FingerprintResult `json:"type"`
	Domain        string            `json:"domain"`
	Url           string            `json:"url"`
	Introspection bool              `json:"introspection,omitempty"`
}

func (o *FingerprintOutput) MarshalJSON() ([]byte, error) {
	// Removes the introspection field from the output if it is disabled
	if !config.Conf.Introspection {
		return json.Marshal(FingerprintOutput{
			Type:   o.Type,
			Domain: o.Domain,
			Url:    o.Url,
		})
	}
	return json.Marshal(FingerprintOutput{
		Type:          o.Type,
		Domain:        o.Domain,
		Url:           o.Url,
		Introspection: o.Introspection,
	})
}
