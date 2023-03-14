package engine

var Ruby = &engine{
	Name: "Ruby",
	Imprints: []imprint{
		{
			Query:   "query @deprecated { __typename }",
			Matcher: inResponseText([]string{"'@deprecated' can't be applied to queries"}),
		},
		{
			Query:   "query @skip { __typename }",
			Matcher: inResponseText([]string{"'@skip' can't be applied to queries (allowed: fields, fragment spreads, inline fragments)"}),
		},
		{
			Query:   "query { __typename @skip }",
			Matcher: inResponseText([]string{"Directive 'skip' is missing required arguments: if"}),
		},
		{
			Query:   "query { __typename {}",
			Matcher: inResponseText([]string{"Parse error on \"}\" (RCURLY)"}),
		},
	},
}
