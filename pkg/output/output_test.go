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
				Domain:        domain,
				Url:           url,
				SchemaStatus:  SchemaStatusOpen,
				Authenticated: false,
				Source:        source,
			},
			&config.Config{
				Introspection:   true,
				FieldSuggestion: false,
			},
			`{"domain":"` + domain + `","url":"` + url + `","authenticated":false,"schema_status":"` + string(SchemaStatusOpen) + `","source":"` + source + `"}`,
		},
		{
			&FingerprintOutput{
				Domain:        domain,
				Url:           url,
				SchemaStatus:  SchemaStatusLeaking,
				Authenticated: true,
				Source:        source,
			},
			&config.Config{
				Introspection:   true,
				FieldSuggestion: true,
			},
			`{"domain":"` + domain + `","url":"` + url + `","authenticated":true,"schema_status":"` + string(SchemaStatusLeaking) + `", "source":"` + source + `"}`,
		},
		{
			&FingerprintOutput{
				Domain:        domain,
				Url:           url,
				SchemaStatus:  SchemaStatusClosed,
				Authenticated: false,
				Source:        source,
			},
			&config.Config{
				Introspection:   false,
				FieldSuggestion: false,
			},
			`{"domain":"` + domain + `","url":"` + url + `","authenticated":false,"schema_status":"` + string(SchemaStatusClosed) + `",` + `"source":"` + source + `"}`,
		},
	}
	for _, table := range tables {
		got, err := marshalOutput(table.output, table.config)
		require.NoError(t, err)
		require.JSONEq(t, table.want, string(got))
	}
}
