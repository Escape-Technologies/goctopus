package engine

var Adriane = &engine{
	Name: "Adriane",
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
