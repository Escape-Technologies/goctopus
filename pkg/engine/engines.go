package engine

import "bytes"

type engine struct {
	Name     string
	Imprints []imprint
}

type imprint struct {
	Query   string
	Matcher matcher
}

// A responseMatcher is a function that takes a response body and returns true if the response matches the engine.
type matcher func(responseBody *[]byte) bool

// inResponseText returns a responseMatcher that checks if the response body contains any of the given strings.
func inResponseText(matches []string) matcher {
	return func(responseBody *[]byte) bool {
		for _, match := range matches {
			if bytes.Contains(*responseBody, []byte(match)) {
				return true
			}
		}
		return false
	}
}

// Order is important here, as the first match will be returned.
// The order has been determined by the usage statistics of the engines. (The higher the usage, the higher the priority.)
var Engines = []*engine{
	Apollo,
	Adriane,
}
