package endpoint

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Escape-Technologies/goctopus/pkg/address"
	"github.com/Escape-Technologies/goctopus/pkg/config"
	out "github.com/Escape-Technologies/goctopus/pkg/output"
)

type mockedEndpointFingerprinter struct {
	openGraphql         bool
	authentifiedGraphql bool
	introspection       bool
	fieldSuggestion     bool
}

func (m *mockedEndpointFingerprinter) IsOpenGraphql() (bool, error) {
	return m.openGraphql, nil
}

func (m *mockedEndpointFingerprinter) IsAuthentifiedGraphql() (bool, error) {
	return m.authentifiedGraphql, nil
}

func (m *mockedEndpointFingerprinter) HasIntrospectionOpen() (bool, error) {
	return m.introspection, nil
}

func (m *mockedEndpointFingerprinter) HasFieldSuggestion() (bool, error) {
	return m.fieldSuggestion, nil
}

func (m *mockedEndpointFingerprinter) Close() {}

func makeMockedEndpointFingerprinter(graphql bool, introspection bool) *mockedEndpointFingerprinter {
	return &mockedEndpointFingerprinter{
		openGraphql:   graphql,
		introspection: introspection,
	}
}

// @todo test field suggestion
func TestFingerprintUrl(t *testing.T) {

	url := &address.Sourced{
		Address: "https://example.com/graphql",
		Source:  "example.com",
	}
	_type := out.ResultOpenGraphql

	table := []struct {
		graphql                    bool
		introspection              bool
		expectedOutput             *out.FingerprintOutput
		expectedErr                error
		configIntrospectionEnabled bool
	}{
		{
			true,
			true,
			&out.FingerprintOutput{
				Introspection: true,
				Url:           url.Address,
				Source:        url.Source,
				Type:          _type,
			},
			nil,
			true,
		},
		{
			true,
			false,
			&out.FingerprintOutput{
				Introspection: false,
				Url:           url.Address,
				Source:        url.Source,
				Type:          _type,
			},
			nil,
			true,
		},
		{
			true,
			false,
			&out.FingerprintOutput{
				Introspection: false,
				Url:           url.Address,
				Source:        url.Source,
				Type:          _type,
			},
			nil,
			false,
		},
		{
			false,
			false,
			nil,
			errors.New("no graphql endpoint found on this route"),
			false,
		},
	}

	for i, test := range table {
		e := makeMockedEndpointFingerprinter(test.graphql, test.introspection)
		config := &config.Config{
			Introspection: test.configIntrospectionEnabled,
		}
		output, err := fingerprintEndpoint(url, e, config)
		if err != nil {
			if err.Error() != test.expectedErr.Error() {
				t.Errorf("expected error %v, got %v", test.expectedErr, err)
			}
		}
		if !reflect.DeepEqual(output, test.expectedOutput) {
			t.Errorf("(case %d) expected output %+v, got %+v", i, test.expectedOutput, output)
		}
		if output == nil && err == nil {
			t.Errorf("Returned no error nor output in case: %d", i)
		}
	}

}
