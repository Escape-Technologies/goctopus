package output

import (
	"testing"

	"github.com/Escape-Technologies/goctopus/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestMarshallOutput(t *testing.T) {
	domain := "example.com"
	url := "https://example.com/graphql"
	source := "example.com"
	tables := []struct {
		output *FingerprintOutput
		config *config.Config
		want   string
	}{
		{
			&FingerprintOutput{
				Domain:          domain,
				Url:             url,
				Introspection:   true,
				FieldSuggestion: false,
				Type:            ResultOpenGraphql,
				Source:          source,
			},
			&config.Config{
				Introspection:   true,
				FieldSuggestion: false,
			},
			`{"domain":"` + domain + `","url":"` + url + `","type":"` + string(ResultOpenGraphql) + `","introspection":true` + `, "source":"` + source + `"}`,
		},
		{
			&FingerprintOutput{
				Domain:          domain,
				Url:             url,
				Introspection:   false,
				FieldSuggestion: true,
				Type:            ResultOpenGraphql,
				Source:          source,
			},
			&config.Config{
				Introspection:   true,
				FieldSuggestion: true,
			},
			`{"domain":"` + domain + `","url":"` + url + `","type":"` + string(ResultOpenGraphql) + `","field_suggestion":true, "introspection":false` + `, "source":"` + source + `"}`,
		},
		{
			&FingerprintOutput{
				Domain:          domain,
				Url:             url,
				Introspection:   false,
				FieldSuggestion: false,
				Type:            ResultAuthentifiedGraphql,
				Source:          source,
			},
			&config.Config{
				Introspection:   false,
				FieldSuggestion: false,
			},
			`{"domain":"` + domain + `","url":"` + url + `","type":"` + string(ResultAuthentifiedGraphql) + `",` + `"source":"` + source + `"}`,
		},
	}
	for _, table := range tables {
		got, err := marshalOutput(table.output, table.config)
		require.NoError(t, err)
		require.JSONEq(t, table.want, string(got))
	}
}
