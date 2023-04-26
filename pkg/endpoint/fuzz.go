package endpoint

import (
	"fmt"

	"github.com/Escape-Technologies/goctopus/pkg/address"
)

func FuzzRoutes(domain *address.Addr) []*address.Addr {
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
	endpoints := []*address.Addr{}
	for _, route := range routes {
		url := fmt.Sprintf("https://%s/%s", domain.Address, route)
		endpoint := domain.Copy()
		endpoint.Address = url
		endpoints = append(endpoints, endpoint)
	}
	return endpoints
}
