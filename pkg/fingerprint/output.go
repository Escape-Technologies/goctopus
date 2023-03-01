package fingerprint

import (
	"encoding/json"

	"github.com/Escape-Technologies/goctopus/internal/config"
)

type OutputType string

const (
	IsGraphql    OutputType = "IS_GRAPHQL"
	MaybeGraphql OutputType = "MAYBE_GRAPHQL"
)

type IsGraphqlOutput struct {
	// The result of the fingerprint if we are sure the domain has a graphql endpoint
	Domain        string
	Url           string
	Introspection bool
}

type MaybeGraphqlOutput struct {
	// The result of the fingerprint if the domain might have a graphql endpoint
	Domain string
	Urls   []string
}

type FingerprintOutput interface {
	GetType() OutputType
	SetDomain(domain string)
}

func (o *IsGraphqlOutput) GetType() OutputType {
	return IsGraphql
}

func (o *MaybeGraphqlOutput) GetType() OutputType {
	return MaybeGraphql
}

func (o *IsGraphqlOutput) SetDomain(domain string) {
	o.Domain = domain
}

func (o *MaybeGraphqlOutput) SetDomain(domain string) {
	o.Domain = domain
}

func (o *IsGraphqlOutput) MarshalJSON() ([]byte, error) {
	type Base struct {
		Type   OutputType `json:"type"`
		Domain string     `json:"domain"`
		Url    string     `json:"url"`
	}

	if !config.Conf.Introspection {
		return json.Marshal(struct {
			Base `json:",inline"`
		}{
			Base: Base{
				Type:   o.GetType(),
				Domain: o.Domain,
				Url:    o.Url,
			},
		})
	}
	return json.Marshal(struct {
		Base          `json:",inline"`
		Introspection bool `json:"introspection"`
	}{
		Base: Base{
			Type:   o.GetType(),
			Domain: o.Domain,
			Url:    o.Url,
		},
		Introspection: o.Introspection,
	})
}

func (o *MaybeGraphqlOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   OutputType `json:"type"`
		Domain string     `json:"domain"`
		Urls   []string   `json:"urls"`
	}{
		Type:   o.GetType(),
		Domain: o.Domain,
		Urls:   o.Urls,
	})
}
