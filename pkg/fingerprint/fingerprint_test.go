package fingerprint

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Escape-Technologies/goctopus/internal/config"
	out "github.com/Escape-Technologies/goctopus/pkg/output"
)

type mockedFingerprinter struct {
	graphql         bool
	introspection   bool
	fieldSuggestion bool
}

func (m *mockedFingerprinter) Graphql() (bool, error) {
	return m.graphql, nil
}

func (m *mockedFingerprinter) Introspection() (bool, error) {
	return m.introspection, nil
}

func (m *mockedFingerprinter) FieldSuggestion() (bool, error) {
	return m.fieldSuggestion, nil
}

func makeMockedFingerprinter(graphql bool, introspection bool) *mockedFingerprinter {
	return &mockedFingerprinter{
		graphql:       graphql,
		introspection: introspection,
	}
}

// @todo test field suggestion
func TestFingerprintUrl(t *testing.T) {

	url := "https://example.com/graphql"
	_type := out.ResultIsGraphql

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
				Url:           url,
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
				Url:           url,
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
				Url:           url,
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
		fp := makeMockedFingerprinter(test.graphql, test.introspection)
		config := &config.Config{
			Introspection: test.configIntrospectionEnabled,
		}
		output, err := FingerprintUrl(url, fp, config)
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
