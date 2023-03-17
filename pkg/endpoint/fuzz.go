package endpoint

import (
	"fmt"

	"github.com/Escape-Technologies/goctopus/pkg/address"
)

func FuzzRoutes(domain *address.Sourced) []*address.Sourced {
	routes := []string{
		"",
		"graphql",
		"graphql/v2",
		"graphql/v1",
		"api",
		"api/graphql",
		"api/v2/graphql",
		"api/v1/graphql",
		"appsync",
		"altair",
		"graph",
	}
	endpoints := []*address.Sourced{}
	for _, route := range routes {
		url := fmt.Sprintf("https://%s/%s", domain.Address, route)
		endpoints = append(endpoints, address.NewSourced(url, domain.Source))
	}
	return endpoints
}
