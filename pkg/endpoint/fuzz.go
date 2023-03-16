package endpoint

import (
	"fmt"
)

func FuzzRoutes(domain string) []string {
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
	endpoints := []string{}
	for _, route := range routes {
		url := fmt.Sprintf("https://%s/%s", domain, route)
		endpoints = append(endpoints, url)
	}
	return endpoints
}
