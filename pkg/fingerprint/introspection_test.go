package fingerprint

import (
	"testing"

	"github.com/Escape-Technologies/goctopus/internal/test/helpers"
	"github.com/valyala/fasthttp"
)

func TestIsValidIntrospectionResponse(t *testing.T) {
	resp := &fasthttp.Response{}
	resp.AppendBody([]byte(`{"data":{"__typename":"Query"}}`))
	tables := []struct {
		resp *fasthttp.Response
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
