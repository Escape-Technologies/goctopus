package output

import (
	"encoding/json"

	"github.com/Escape-Technologies/goctopus/internal/utils"
	"github.com/Escape-Technologies/goctopus/pkg/config"
)

type FingerprintResult string

const (
	ResultOpenGraphql         FingerprintResult = "OPEN_GRAPHQL"
	ResultAuthentifiedGraphql FingerprintResult = "AUTHENTIFIED_GRAPHQL"
)

type FingerprintOutput struct {
	Type            FingerprintResult `json:"type"`
	Domain          string            `json:"domain"`
	Url             string            `json:"url"`
	Introspection   bool              `json:"introspection"`
	FieldSuggestion bool              `json:"field_suggestion"`
	Source          string            `json:"source"`   // the original address used to fingerprint the endpoint
	Metadata        map[string]string `json:"metadata"` // optional metadata
}

func (o *FingerprintOutput) MarshalJSON() ([]byte, error) {
	return marshalOutput(o, config.Get())
}

// This is separated from the above function to decouple the output from the config
func marshalOutput(o *FingerprintOutput, c *config.Config) ([]byte, error) {
	// Removes the introspection field from the output if it is disabled

	// transform the output to a map
	outputMap := make(map[string]interface{})
	// this is need to avoid an infinite recursion when marshaling the output
	type alias FingerprintOutput
	outputBytes, _ := json.Marshal((*alias)(o))
	if err := json.Unmarshal(outputBytes, &outputMap); err != nil {
		panic(err)
	}

	// remove the introspection field if it is disabled
	if !c.Introspection {
		delete(outputMap, "introspection")
	}

	// if introspection is enabled, remove the field_suggestion field
	if !c.FieldSuggestion || o.Introspection {
		delete(outputMap, "field_suggestion")
	}

	// when scanning from an url, the domain is not set so we infer it from the url
	if o.Domain == "" && o.Url != "" {
		outputMap["domain"] = utils.DomainFromUrl(o.Url)
	}

	if o.Metadata == nil || len(o.Metadata) == 0 {
		delete(outputMap, "metadata")
	}

	return json.Marshal(outputMap)
}
