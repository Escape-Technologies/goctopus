package output

import (
	"encoding/json"

	"github.com/Escape-Technologies/goctopus/internal/config"
	"github.com/Escape-Technologies/goctopus/internal/utils"
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
	Introspection   bool              `json:"introspection"`
	FieldSuggestion bool              `json:"field_suggestion"`
}

func (o *FingerprintOutput) MarshalJSON() ([]byte, error) {
	return marshalOutput(o, config.Conf)
}

// This is separated from the above function to decouple the output from the config
func marshalOutput(o *FingerprintOutput, c *config.Config) ([]byte, error) {
	// Removes the introspection field from the output if it is disabled

	// transform the output to a map
	outputMap := make(map[string]interface{})
	// this is need to avoid an infinite recursion
	type alias FingerprintOutput
	outputBytes, _ := json.Marshal((*alias)(o))
	json.Unmarshal(outputBytes, &outputMap)

	// remove the introspection field if it is disabled
	if !c.Introspection {
		delete(outputMap, "introspection")
	}

	if !c.FieldSuggestion {
		delete(outputMap, "field_suggestion")
	}

	// when scanning from an url, the domain is not set so we infer it from the url
	if o.Domain == "" && o.Url != "" {
		outputMap["domain"] = utils.DomainFromUrl(o.Url)
	}

	return json.Marshal(outputMap)
}
