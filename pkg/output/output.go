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
	Type            FingerprintResult `json:"type"`
	Domain          string            `json:"domain"`
	Url             string            `json:"url"`
	Introspection   bool              `json:"introspection,omitempty"`
	FieldSuggestion bool              `json:"field_suggestion,omitempty"`
}

// @todo: make this work -> segfaults
// func (o *FingerprintOutput) MarshalJSON() ([]byte, error) {
// 	// Removes the introspection field from the output if it is disabled

// 	// transform the output to a map
// 	var outputMap map[string]interface{}
// 	outputBytes, _ := json.Marshal(o)
// 	json.Unmarshal(outputBytes, &outputMap)

// 	// remove the introspection field if it is disabled
// 	if !config.Conf.Introspection {
// 		delete(outputMap, "introspection")
// 	}

// 	return json.Marshal(outputMap)
// }

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
