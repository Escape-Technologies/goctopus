package engine

import (
	"bytes"
	"encoding/json"
)

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

// inSection returns a responseMatcher that checks if the response body contains any of the given strings in the given section.
func inSection(section string, matches []string) matcher {
	return func(responseBody *[]byte) bool {
		var reponseBody map[string]interface{}
		json.Unmarshal(*responseBody, &reponseBody)
		content, err := json.Marshal(reponseBody[section])
		if err != nil {
			return false
		}
		for _, match := range matches {
			if bytes.Contains(content, []byte(match)) {
				return true
			}
		}
		return false
	}
}

// hasJsonKey returns a responseMatcher that checks if the response body contains the given key.
func hasJsonKey(key string) matcher {
	return func(responseBody *[]byte) bool {
		var reponseBody map[string]interface{}
		json.Unmarshal(*responseBody, &reponseBody)
		_, ok := reponseBody[key]
		return ok
	}
}

// Order is important here, as the first match will be returned.
// The order has been determined by the usage statistics of the engines. (The higher the usage, the higher the priority.)
var Engines = []*engine{
	Apollo,
	AWSAppSync,
	GraphQLGo,
	Ruby,
	GraphQLPHP,
	Graphene,
	Adriane,
	GraphQLGopherGo,
}
