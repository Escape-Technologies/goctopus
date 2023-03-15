package fingerprint

import (
	"testing"

	"github.com/Escape-Technologies/goctopus/internal/http"
	"github.com/Escape-Technologies/goctopus/internal/test/helpers"
)

func TestIsAuthentifiedGraphqlResponse(t *testing.T) {
	tables := []struct {
		resp *http.Response
		want bool
	}{
		{
			helpers.MockHttpResponse(
				200,
				`{"data":{"__typename":"Query"}}`,
			),
			false,
		},
		{
			helpers.MockHttpResponse(
				200,
				`{"errors":[]}`,
			),
			false,
		},
		{
			helpers.MockHttpResponse(
				200,
				`{
					"errors": [{
						"message": "You must be logged in to perform this action",
						"locations": [{
							"line": 5,
							"column": 2
						}]
					}]
				}`,
			),
			true,
		},
		{
			helpers.MockHttpResponse(
				403,
				`{
					"errors": [{
						"message": "Unauthorized"
					}]
				}`,
			),
			true,
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
		got := isAuthentifiedGraphqlResponse(table.resp)
		if got != table.want {
			t.Errorf("isAuthentifiedGraphqlResponse() was incorrect for %+v, got: %v, want: %v.", table.resp, got, table.want)
		}
	}

}

func TestIsValidGraphqlResponse(t *testing.T) {
	tables := []struct {
		resp *http.Response
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
		got := isOpenGraphqlResponse(table.resp)
		if got != table.want {
			t.Errorf("isOpenGraphqlResponse() was incorrect for %+v, got: %v, want: %v.", table.resp, got, table.want)
		}
	}

}
