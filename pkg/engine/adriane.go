package engine

var Adriane = &engine{
	Name: "adriane",
	Imprints: []imprint{
		{
			Query:   "",
			Matcher: inResponseText([]string{"The query must be a string."}),
		},
		{
			Query:   "query { __typename @abc }",
			Matcher: inResponseText([]string{"Unknown directive '@abc'."}),
		},
	},
}
