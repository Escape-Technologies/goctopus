package fingerprint

import (
	"testing"

	"github.com/Escape-Technologies/goctopus/internal/http"
	"github.com/Escape-Technologies/goctopus/internal/test/helpers"
)

func TestMatchFieldSuggestionRegex(t *testing.T) {
	tables := []struct {
		message string
		want    bool
	}{
		{
			"Cannot query field \"tes\" on type \"Query\". Did you mean \"test\"?",
			true,
		},
		{
			"ERROR at Object.Field...",
			false,
		},
	}

	for _, table := range tables {
		got := MatchFieldSuggestionRegex(table.message)
		if got != table.want {
			t.Errorf("got: '%v', expected: '%v' for message '%v'", got, table.want, table.message)
		}
	}
}

func TestIsSuggestionResponse(t *testing.T) {
	tables := []struct {
		resp *http.Response
		want bool
	}{
		{
			helpers.MockHttpResponse(
				200,
				`{
					"errors": 
						[
							{ 
								"message": "Cannot query field \"tes\" on type \"Query\". Did you mean \"test\"?"
							}
						]
					}`,
			),
			true,
		},
		{
			helpers.MockHttpResponse(
				400,
				`{
					"errors": 
						[
							{ 
								"message": "ERROR at Object.Field..."
							},
							{ 
								"message": "Cannot query field \"tes\" on type \"Query\". Did you mean \"test\"?"
							}
						]
					}`,
			),
			true,
		},
		{
			helpers.MockHttpResponse(
				200,
				`ERROR at Object.Field...`,
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
		got := IsSuggestionResponse(table.resp)
		if got != table.want {
			t.Errorf("IsSuggestionResponse() was incorrect for %+v, got: %v, want: %v.", table.resp, got, table.want)
		}
	}
}
