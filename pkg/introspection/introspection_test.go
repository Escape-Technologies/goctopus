package introspection

import (
	"testing"

	"github.com/Escape-Technologies/goctopus/internal/test/helpers"
	"github.com/Escape-Technologies/goctopus/pkg/http"
)

func TestIsValidIntrospectionResponse(t *testing.T) {
	tables := []struct {
		resp *http.Response
		want bool
	}{
		{
			helpers.MockHttpResponse(
				200,
				`{"data":{"__schema":{"queryType":{"name":"Query"}}}}`,
			),
			true,
		},
		{
			helpers.MockHttpResponse(
				200,
				`{"data":{"__schema":{"queryType":{"name":""}}}}`,
			),
			false,
		},
		{
			helpers.MockHttpResponse(
				200,
				`{"data":null, "errors": [{"message": "Error"}]}`,
			),
			false,
		},
		{
			helpers.MockHttpResponse(
				200,
				"Error",
			),
			false,
		},
		{
			helpers.MockHttpResponse(
				500,
				"",
			),
			false,
		},
	}

	for _, table := range tables {
		got := IsValidIntrospectionResponse(table.resp)
		if got != table.want {
			t.Errorf("IsValidGraphqlResponse() was incorrect for %+v, got: %v, want: %v.", table.resp, got, table.want)
		}
	}
}
