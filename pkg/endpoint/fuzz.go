package endpoint

import (
	"fmt"

	"github.com/Escape-Technologies/goctopus/pkg/address"
)

/** @todo: remove less probable routes one script analysis is done, or add a slow mode or dictionnary choice */
func FuzzRoutes(domain *address.Addr) []*address.Addr {
	routes := []string{
		"graphql",
		"api/graphql",
		"",
		"api",
		"graphql/v2",
		"v1/graphql",
		"graphql/v1",
		"api/v2/graphql",
		"graphql/console",
		"api/v1/graphql",
		"graph",
		"dev/graphql",
		"v1",
		"appsync",
		"altair",
		"graph/api",
		"v2/graphql",
		"v1/graphql-public",
		"v2/graphql-public",
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
