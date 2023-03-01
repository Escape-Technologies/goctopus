package fingerprint

import (
	"testing"

	"github.com/Escape-Technologies/goctopus/test/helpers"
	"github.com/valyala/fasthttp"
)

func TestIsValidGraphqlResponse(t *testing.T) {
	resp := &fasthttp.Response{}
	resp.AppendBody([]byte(`{"data":{"__typename":"Query"}}`))
	tables := []struct {
		resp *fasthttp.Response
		want bool
	}{
		{
			helpers.MockHttpResponse(
				200,
				`{"data":{"__typename":"Query"}}`,
			),
			true,
		},
		{
			helpers.MockHttpResponse(
				200,
				`{"data":{"__typename":""}}`,
			),
			false,
		},
		{
			helpers.MockHttpResponse(
				200,
				`{"data":{"__typename":null}}`,
			),
			false,
		},
		{
			helpers.MockHttpResponse(
				200,
				`{"data":{}}`,
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
		got := IsValidGraphqlResponse(table.resp)
		if got != table.want {
			t.Errorf("IsValidGraphqlResponse() was incorrect for %+v, got: %v, want: %v.", table.resp, got, table.want)
		}
	}

}
